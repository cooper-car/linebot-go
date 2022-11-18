package main

import (
	"fmt"
	// "io/ioutil"
	"log"
    "strconv"
	"net/http"
    "github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
)

var bot *linebot.Client
var collection *mongo.Collection

type Message struct {
	Type       string  `json:"id"`
	Text       string  `json:"text"`
	ReplyToken string  `json:"reply_token"`
}

func init() {
    // DB連線
	clientOptions := options.Client().ApplyURI("mongodb://root:123456@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
    collection = client.Database("linebot").Collection("message")
	fmt.Println("Connected to MongoDB!")
}


func main() {
    var err error
    // 取設定檔資料
    viper.SetConfigName("app")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    error := viper.ReadInConfig()
    if error != nil {
        panic(error)
    }

    // 創建line bot
    bot, err = linebot.New(viper.GetString("application.ChannelSecret"), viper.GetString("application.ChannelAccessToken"))
    if err != nil {
        fmt.Println("Bot:", bot, " err:", err)
	}

    // router
    router := gin.Default()
    router.GET("/quota", getQuota)
    router.POST("/pushMessage", pushMessage)
    router.POST("/callback", callback)
    router.GET("/queryMessage", queryMessage)
    
    // port 
    addr := fmt.Sprintf(":%s", viper.GetString("application.port"))
    router.Run(addr)
   
}

// 發送通知
func pushMessage(c *gin.Context) {
    inputMsg := c.PostForm("message")
    message := linebot.NewTextMessage(inputMsg)
    _, err := bot.PushMessage(viper.GetString("application.userId"), message).Do()
    if err != nil {
        log.Print(err)
    }

    c.JSON(200, gin.H{
        "data": message,
    })
}

// user接收發送通知
func callback(c *gin.Context) {
    events, err := bot.ParseRequest(c.Request)
    if err != nil {
        if err == linebot.ErrInvalidSignature {
            log.Print(err)
        }
        return
    }

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					fmt.Println("Quota err:", err)
				}

                msg := Message{string(event.Type), message.Text, event.ReplyToken}
                insertResult, err := collection.InsertOne(context.TODO(), msg)
                if err != nil {
                    log.Fatal(err)
                }

                fmt.Println("Insert to mongodb: ", insertResult)
        
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("content:"+message.Text+" , \nremain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					fmt.Print(err)
				}
			}
		}
	}
}

// 查詢mongodb裡面line bot接收的訊息資料
func queryMessage(c *gin.Context){
    var results []Message
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Message
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

    cur.Close(context.TODO())

    c.JSON(200, gin.H{
        "data": results,
    })
}

// 測試：剩餘訊息發送量
func getQuota(c *gin.Context) {
    quota, err := bot.GetMessageQuota().Do()
    if err != nil {
        fmt.Println("Quota err:", err)
    }

    c.JSON(http.StatusOK, gin.H{
		"quota":   quota,
	})
}


