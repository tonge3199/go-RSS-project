package main

import (
	"context"
	"encoding/xml"
	"github.com/tonge3199/go-RSS-project/internal/database"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// startScraping 周期性抓取订阅源数据，使用指定数量的 goroutine 并行处理
// 参数:
//   - db: 数据库查询接口，用于获取待抓取订阅源和存储结果
//   - concurrency: 最大并发处理数，控制同时处理的订阅源数量
//   - timeBetweenRequest: 批次之间的等待时间，控制抓取频率
//
// 函数行为:
// 1. 按照指定时间间隔周期性触发抓取任务
// 2. 每次获取最多 concurrency 个待处理订阅源
// 3. 使用单独的 goroutine 并行处理每个订阅源
// 4. 等待当前批次全部完成后进入下一个周期
// 5. 发生数据库错误时记录日志并继续运行
// 6. 函数会阻塞当前 goroutine 无限运行，通常需要在单独 goroutine 中调用
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)

	// 永久运行循环，使用 ticker 控制节奏
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue // 跳过当前周期，等待下次重试
		}
		log.Printf("Fetching %d feeds on %v goroutines", len(feeds), concurrency)

		// 2. 并发处理所有订阅源
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

// scrapeFeed 执行单个订阅源的抓取和处理流程
//
// 参数:
//   - db: 数据库操作接口，用于更新订阅源状态
//   - wg: 等待组对象，用于同步并发任务
//   - feed: 要抓取的订阅源元数据
//
// 处理流程:
// 1. 立即标记订阅源为已抓取状态（无论后续是否成功）
// 2. 尝试获取远程订阅源内容
// 3. 遍历所有找到的内容条目
// 4. 记录最终抓取结果统计
//
// 注意:
// - 使用 wg.Done() 确保等待组计数器递减
// - 先标记后抓取的设计适用于避免重复抓取失败订阅源
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // 确保递减等待组计数器

	// 标记订阅源为已抓取状态
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	// 获取远程订阅源内容
	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	// 处理每个内容条目
	for _, item := range feedData.Channel.Item {
		log.Println("Found item", item.Title, item.Link)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

// RSSFeed 表示整个RSS文档的根结构
// XML对应结构示例：
// <rss>
//
//	<channel>...</channel>
//
// </rss>
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem 表示单个文章条目
// XML对应结构示例：
// <item>
//
//	<title>...</title>
//	<link>...</link>
//	<description>...</description>
//	<pubDate>...</pubDate>
//
// </item>
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}
