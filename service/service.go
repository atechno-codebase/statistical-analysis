package service

import (
	"context"
	"errors"
	"fmt"
	"statistical-analysis/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var databaseName string

func Init(mongoUrl, dbName string) {
	databaseName = dbName
	database.Init(mongoUrl)
}

func CalculateStats(uid string) (interface{}, error) {

	statData, err := GetStatsFromDB(uid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return statData, err
}

func GetNode(uid string) (bson.M, error) {
	node, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		collection := client.Database(databaseName).Collection("nodes")

		result, err := collection.Find(context.Background(), bson.M{"uid": uid})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		nodes := []bson.M{}
		for result.Next(context.TODO()) {
			node := bson.M{}
			result.Decode(&node)
			nodes = append(nodes, node)
		}
		if len(nodes) > 0 {
			return nodes[0], err
		}
		return nil, errors.New("Node not found")
	})
	return node.(bson.M), err
}

func GetStatsFromDB(uid string) (interface{}, error) {
	node, err := GetNode(uid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	filter := bson.M{}
	if node["isTemperature"].(bool) {
		filter["temperature"] = bson.M{"$exists": true}
	}
	if node["isHumidity"].(bool) {
		filter["humidity"] = bson.M{"$exists": true}
	}
	if node["isCO2"].(bool) {
		filter["co2"] = bson.M{"$exists": true}
	}

	query := bson.A{
		bson.M{
			"$match": bson.M{
				"uid": uid,
			},
		},
		bson.M{
			"$match": filter,
		},
		bson.M{
			"$group": bson.M{
				"_id":              "$uid",
				"temperatureSigma": bson.M{"$stdDevSamp": "$temperature"},
				"humiditySigma":    bson.M{"$stdDevSamp": "$humidity"},
				"co2Sigma":         bson.M{"$stdDevSamp": "$co2"},
				"temperatureMean":  bson.M{"$avg": "$temperature"},
				"humidityMean":     bson.M{"$avg": "$humidity"},
				"co2Mean":          bson.M{"$avg": "$co2"},
			},
		},
	}

	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		collection := client.Database(databaseName).Collection("readings")

		result, err := collection.Aggregate(context.Background(),
			query,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		var final = []bson.M{}
		for result.Next(context.TODO()) {
			var iResult bson.M
			if err := result.Decode(&iResult); err != nil {
				fmt.Println(err)
				return nil, err
			}
			final = append(final, iResult)
		}

		return final, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := res.([]bson.M)

	faulty_readings, err := GetFaultyReadings(uid, node)
	if err != nil {
		fmt.Println(err)
		response[0]["faulty_readings"] = 0
		return response[0], nil
	}
	response[0]["faulty_readings"] = faulty_readings.(bson.M)["faulty_readings"]
	delete(response[0], "_id")
	response[0]["uid"] = uid

	return response[0], nil
}

func GetFaultyReadings(uid string, node bson.M) (interface{}, error) {
	filter := bson.A{}

	if node["isTemperature"].(bool) {
		filter = append(filter, bson.M{
			"$or": bson.A{
				bson.M{
					"temperature": bson.M{"$gt": node["temperatureRange"].(bson.M)["max"]},
				},
				bson.M{
					"temperature": bson.M{"$lt": node["temperatureRange"].(bson.M)["min"]},
				},
			},
		})
	}
	if node["isHumidity"].(bool) {
		filter = append(filter, bson.M{
			"$or": bson.A{
				bson.M{
					"humidity": bson.M{"$gt": node["humidityRange"].(bson.M)["max"]},
				},
				bson.M{
					"humidity": bson.M{"$lt": node["humidityRange"].(bson.M)["min"]},
				},
			},
		})
	}
	if node["isCO2"].(bool) {
		filter = append(filter, bson.M{
			"$or": bson.A{
				bson.M{
					"co2": bson.M{"$gt": node["co2Range"].(bson.M)["max"]},
				},
				bson.M{
					"co2": bson.M{"$lt": node["co2Range"].(bson.M)["min"]},
				},
			},
		})
	}

	mainQuery := bson.A{
		bson.M{
			"$match": bson.M{
				"uid": uid,
			},
		},
		bson.M{
			"$match": bson.M{
				"$or": filter,
			},
		},
		bson.M{
			"$count": "faulty_readings",
		},
	}

	fmt.Println("main query: ", mainQuery)

	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		collection := client.Database(databaseName).Collection("readings")
		result, err := collection.Aggregate(context.Background(), mainQuery)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		var entry interface{}
		result.Next(context.Background())
		err = result.Decode(&entry)
		if err != nil {
			fmt.Println("err: ", err)
			return nil, err
		}
		fmt.Println("response:", entry)
		return entry, nil

		// var entries []bson.M
		// for result.Next(context.TODO()) {
		// 	var entry bson.M
		// 	err := result.Decode(&entry)
		// 	fmt.Println("entry", entry)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return nil, err
		// 	}
		// 	entries = append(entries, entry)
		// }
		// fmt.Println("entries", entries)
		// if len(entries) > 0 {
		// 	return entries[0], nil
		// }
		// return nil, errors.New("No entries found for uid")

	})
	return res, err
}
