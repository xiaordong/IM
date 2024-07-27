package service

import (
	"chat/conf"
	"chat/model"
	"chat/model/ws"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"time"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

// 关于信息的寻找及排序

func InsertMsg(database string, id string, content string, read uint, expire int64) (err error) {
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err = collection.InsertOne(context.TODO(), comment)
	return
}

func FindMany(database, sendId, id string, Time int64, pageSize int) (results []ws.Result, err error) {
	var resultMe []ws.Trainer
	var resultYou []ws.Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)

	// Perform query on sendIdCollection
	sendIdTimeCursor, err := sendIdCollection.Find(context.TODO(), bson.D{
		{"startTime", bson.D{{"$gt", Time}}},
	},
		options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)))
	defer sendIdTimeCursor.Close(context.TODO())
	err = sendIdTimeCursor.All(context.TODO(), &resultYou)
	idTimeCursor, err := idCollection.Find(context.TODO(), bson.D{
		{"startTime", bson.D{{"$gt", Time}}},
	},
		options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		return nil, err
	}
	defer idTimeCursor.Close(context.TODO())
	err = idTimeCursor.All(context.TODO(), &resultMe)
	results, _ = AppendAndSort(resultMe, resultYou)
	return results, nil
}

func FirstFindMsg(database string, sendId string, id string) (results []ws.Result, err error) {
	// 首次查询(把对方发来的所有未读都取出来)
	var resultsMe []ws.Trainer
	var resultsYou []ws.Trainer
	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	filter := bson.M{"read": bson.M{
		"$all": []uint{0},
	}}
	sendIdCursor, err := sendIdCollection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{
		"startTime", 1}}), options.Find().SetLimit(1))
	if sendIdCursor == nil {
		return
	}
	var unReads []ws.Trainer
	err = sendIdCursor.All(context.TODO(), &unReads)
	if err != nil {
		log.Println("sendIdCursor err", err)
	}
	if len(unReads) > 0 {
		timeFilter := bson.M{
			"startTime": bson.M{
				"$gte": unReads[0].StartTime,
			},
		}
		sendIdTimeCursor, _ := sendIdCollection.Find(context.TODO(), timeFilter)
		idTimeCursor, _ := idCollection.Find(context.TODO(), timeFilter)
		err = sendIdTimeCursor.All(context.TODO(), &resultsYou)
		err = idTimeCursor.All(context.TODO(), &resultsMe)
		results, err = AppendAndSort(resultsMe, resultsYou)
	} else {
		results, err = FindMany(database, sendId, id, 9999999999, 10)
	}
	overTimeFilter := bson.D{
		{"$and", bson.A{
			bson.D{{"endTime", bson.M{"&lt": time.Now().Unix()}}},
			bson.D{{"read", bson.M{"$eq": 1}}},
		}},
	}
	_, _ = sendIdCollection.DeleteMany(context.TODO(), overTimeFilter)
	_, _ = idCollection.DeleteMany(context.TODO(), overTimeFilter)
	// 将所有的维度设置为已读
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"read": 1},
	})
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"&set": bson.M{"ebdTime": time.Now().Unix() + int64(3*month)},
	})
	return
}

func AppendAndSort(resultsMe, resultsYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultsMe {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, r := range resultsYou {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	// 最后进行排序
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil
}

// 关于用户的查询

func FindUserByInfo(info string) (results model.User, err error) {
	res := model.DB.Model(&model.User{}).Where("user_name = ? OR id = ?", info, info).Find(&results)
	if res.RowsAffected == 0 || res.Error != nil {
		err = res.Error
		log.Println("FindUserByInfo err:", res.Error.Error())
		return model.User{}, err
	}
	return results, nil
}
