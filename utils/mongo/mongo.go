package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

var Client *mongo.Client
var database *mongo.Database

func init() {
	logs.Debug("mongo init")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username, _ := beego.AppConfig.String("mongo.user")
	passwd, _ := beego.AppConfig.String("mongo.pass")
	host, _ := beego.AppConfig.String("mongo.host")
	port, _ := beego.AppConfig.String("mongo.port")
	db, _ := beego.AppConfig.String("mongo.db")
	logs.Debug("mongo info", username, passwd, host, port, db)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port)

	//client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(options.Credential{
	//	Username: username,
	//	Password: passwd,
	//}))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	if err != nil {
		logs.Debug("mongo db Connect err ", err.Error())
		panic(err.Error())
	}
	//err = client.Ping(ctx, readpref.Primary())
	//if err != nil {
	//	logs.Debug("mongo db Connect err ", err.Error())
	//	panic(err.Error())
	//}
	Client = client
	database = client.Database(db)
}

func Insert(i interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := reflect.TypeOf(i)
	// 检查是否是指针，如果是指针则获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 打印结构体的名称
	logs.Debug("结构体名称 = ", t.Name())
	collection := database.Collection(t.Name())

	res, err := collection.InsertOne(ctx, i)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func InsertMany(i []interface{}) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collectionName := ""
	for _, value := range i {
		t := reflect.TypeOf(value)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		logs.Debug("结构体名称 = ", t.Name())
		collectionName = t.Name()
	}

	collection := database.Collection(collectionName)

	res, err := collection.InsertMany(ctx, i)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, err
}

func Update(i interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	// 检查是否是指针，如果是指针则获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	// 打印结构体的名称
	logs.Debug("结构体名称 =", t.Name())
	collection := database.Collection(t.Name())

	idField := v.FieldByName("Id")
	if !(idField.IsValid() && idField.CanInterface()) {
		logs.Debug("结构体中没有 ID 字段或无法访问")
		return nil, errors.New("_id field must set")
	}
	filter := bson.M{"_id": idField.Interface()}
	bsonData, err := bson.Marshal(i)
	if err != nil {
		logs.Debug("BSON 编码失败", err)
		return nil, err
	}

	var bsonDoc bson.D
	err = bson.Unmarshal(bsonData, &bsonDoc)
	if err != nil {
		logs.Debug("BSON 反序列化失败", err)
		return nil, err
	}
	updateData := bson.M{}
	for _, elem := range bsonDoc {
		updateData[elem.Key] = elem.Value
	}

	update := bson.M{"$set": updateData}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Debug("UpdateOne", err)
	}
	return res.UpsertedID, err
}

func Read(i interface{}, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := reflect.TypeOf(i)
	// 检查是否是指针，如果是指针则获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 打印结构体的名称
	logs.Debug("结构体名称 =", t.Name())
	collection := database.Collection(t.Name())

	finalFilter := bson.M{}
	if filter != nil {
		v := reflect.ValueOf(filter)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			// 遍历结构体字段
			for i := 0; i < v.NumField(); i++ {
				fieldName := v.Type().Field(i).Name
				fieldValue := v.Field(i).Interface()
				finalFilter[fieldName] = fieldValue
				logs.Debug("fieldName: fieldValue\n", fieldName, fieldValue)
			}
		case reflect.Map:
			// 遍历 map 键和值
			for _, key := range v.MapKeys() {
				value := v.MapIndex(key).Interface()
				keyStr := fmt.Sprintf("%v", key.Interface())
				logs.Debug("keyStr: value\n", keyStr, value)
				finalFilter[keyStr] = value
			}
		default:
			logs.Debug("Unsupported type")
			return errors.New("filter type is Unsupported")
		}
	}
	err := collection.FindOne(ctx, finalFilter).Decode(i)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	return err
}

func Query(collectionName string, page int, size int, filter interface{}) (int64, []bson.D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.Collection(collectionName)

	finalFilter := bson.M{}
	if filter != nil {
		v := reflect.ValueOf(filter)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			// 遍历结构体字段
			for i := 0; i < v.NumField(); i++ {
				fieldName := v.Type().Field(i).Name
				fieldValue := v.Field(i).Interface()
				if fieldName == "created_time" {
					t := reflect.TypeOf(fieldValue)
					logs.Debug("created_time type = ", t)
					interfaceSlice, ok := fieldValue.([]interface{})
					if !ok {
						return 0, nil, errors.New("created_time must be a Array")
					}
					if len(interfaceSlice) < 2 {
						return 0, nil, errors.New("created_time len must = 2")
					}
					var createdTime []time.Time
					for _, v := range interfaceSlice {
						if str, ok := v.(string); ok {
							startTime, err := time.Parse(time.DateTime, str)
							if err != nil {
								logs.Debug(err)
							}
							createdTime = append(createdTime, startTime)
						} else {
							fmt.Printf("类型断言失败，遇到非 string 类型的值: %v\n", v)
						}
					}

					finalFilter[fieldName] = bson.M{
						"$gte": createdTime[0],
						"$lt":  createdTime[1],
					}
					continue
				}
				finalFilter[fieldName] = fieldValue

				logs.Debug("fieldName: fieldValue\n", fieldName, fieldValue)
			}
		case reflect.Map:
			// 遍历 map 键和值
			for _, key := range v.MapKeys() {
				value := v.MapIndex(key).Interface()
				keyStr := fmt.Sprintf("%v", key.Interface())
				logs.Debug("keyStr: value\n", keyStr, value)
				if keyStr == "created_time" {
					t := reflect.TypeOf(value)
					logs.Debug("created_time type = ", t)
					interfaceSlice, ok := value.([]interface{})
					if !ok {
						return 0, nil, errors.New("created_time must be a Array")
					}
					if len(interfaceSlice) < 2 {
						return 0, nil, errors.New("created_time len must = 2")
					}
					var createdTime []time.Time
					for _, v := range interfaceSlice {
						if str, ok := v.(string); ok {
							startTime, err := time.Parse(time.DateTime, str)
							if err != nil {
								logs.Debug(err)
							}
							createdTime = append(createdTime, startTime)
						} else {
							fmt.Printf("类型断言失败，遇到非 string 类型的值: %v\n", v)
						}
					}

					finalFilter[keyStr] = bson.M{
						"$gte": createdTime[0],
						"$lt":  createdTime[1],
					}
					continue
				}
				finalFilter[keyStr] = value
			}
		default:
			logs.Debug("Unsupported type")
			return 0, nil, errors.New("filter type is Unsupported")
		}
	}

	total, err := collection.CountDocuments(ctx, finalFilter)
	if err != nil {
		logs.Debug("Failed to count documents: %v", err)
	}
	fmt.Printf("Total documents count: %d\n", total)

	skip := (page - 1) * size
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(size))

	logs.Debug("finalFilter ", finalFilter)
	cur, err := collection.Find(ctx, finalFilter, findOptions)
	if err != nil {
		logs.Debug("Find cur", err)
		return 0, nil, err
	}
	var back []bson.D
	for cur.Next(ctx) {
		var result bson.D
		err = cur.Decode(&result)
		if err != nil {
			logs.Debug("cur.Next", err)
			return 0, nil, err
		}
		back = append(back, result)
		// do something with result....
	}
	if err = cur.Err(); err != nil {
		logs.Debug("cur.Err()", err)
		return 0, nil, err
	}
	//jsonData, err := json.Marshal(back)
	//if err != nil {
	//	fmt.Println("JSON marshaling failed:", err)
	//	return nil, err
	//}
	//logs.Debug("jsonData", jsonData)
	return total, back, nil
}
