package mongo

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"tinyUrl/config"
	"tinyUrl/config/log"
)

type B bson.M
type Regex bson.RegEx

var (
	s          *mgo.Session
	defaultDb  string
	SessionMap = make(map[string]*mgo.Session)
)

func getSession() *mgo.Session {
	if s == nil {
		err := Init()
		if err != nil {

		}
	}
	return s.Clone()
}

func Init() error {
	var (
		url      string
		user     string
		password string
		host     string
	)

	user = config.Base.Mongo.User
	password = config.Base.Mongo.Password
	host = config.Base.Mongo.Host
	url = fmt.Sprintf("mongodb://%s:%s@%s/admin", user, password, host)

	session, err := mgo.Dial(url)
	if err != nil {
		log.GetLogger().Errorf("Init MongoDB Error: %v", err)
		return err
	}

	session.SetMode(mgo.Monotonic, false)

	size := config.Base.Mongo.PoolSize
	session.SetPoolLimit(size)

	s = session
	defaultDb = config.Base.Mongo.DbName

	log.GetLogger().Infof("Init MongoDB Success: %s", url)
	return nil
}

func Insert(c string, docs interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Insert(docs)
}

func DbInsert(db, c string, docs interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Insert(docs)
}

func DbInsertMany(db, c string, docs []interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Insert(docs...)
}

func Find(c string, query B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Find(query).All(result)
}

func DbFind(db, c string, query B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).All(result)
}

func DbGetCount(db, c string, query B) (int, error) {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Count()
}

func FindPage(c string, query B, skip int, limit int, result interface{}) error {
	session := getSession()
	defer session.Close()

	if limit != 0 {
		return session.DB(defaultDb).C(c).Find(query).Limit(limit).Skip(skip).All(result)
	} else {
		return session.DB(defaultDb).C(c).Find(query).All(result)
	}
}

func DbFindPage(dbName, c string, query B, skip int, limit int, result interface{}) error {
	session := getSession()
	defer session.Close()

	if limit != 0 {
		return session.DB(dbName).C(c).Find(query).Limit(limit).Skip(skip).All(result)
	} else {
		return session.DB(dbName).C(c).Find(query).All(result)
	}
}

func FindPageSort(c string, query B, skip int, limit int, sort string, result interface{}) error {
	session := getSession()
	defer session.Close()

	if limit != 0 {
		return session.DB(defaultDb).C(c).Find(query).Limit(limit).Sort(sort).Skip(skip).All(result)
	} else {
		return session.DB(defaultDb).C(c).Find(query).Sort(sort).All(result)
	}
}

func DbFindSort(db string, c string, query B, sort string, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Sort(sort).All(result)
}

func DbFindSortDouble(db string, c string, query B, sort1 string, sort2 string, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Sort(sort1, sort2).All(result)
}

func DbFindPageSortDouble(db string, c string, query B, skip int, limit int, sort1 string, sort2 string, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Sort(sort1, sort2).Limit(limit).Skip(skip).All(result)
}

func DBFindPageSort(db string, c string, query B, skip int, limit int, sort string, result interface{}) error {
	session := getSession()
	defer session.Close()
	if sort == "" {
		return errors.New("sort is not exists")
	}
	if limit != 0 {
		return session.DB(db).C(c).Find(query).Sort(sort).Limit(limit).Skip(skip).All(result)
	} else {
		return session.DB(db).C(c).Find(query).Sort(sort).All(result)
	}
}

func DbFindByFields(db string, c string, query B, field B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Select(field).All(result)
}

func FindByFields(c string, query B, field B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Find(query).Select(field).All(result)
}

func FindOne(c string, query B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Find(query).One(result)
}

func DbFindOne(db string, c string, query B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).One(result)
}

func DbFindById(db string, c string, id interface{}, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).FindId(id).One(result)
}

func FindById(c string, id interface{}, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).FindId(id).One(result)
}

func UpdateOne(c string, query B, set B) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Update(query, set)
}

func DbUpdateOne(db, c string, query B, set B) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Update(query, set)
}

func DBUpdateById(DB string, c string, id interface{}, set B) error {
	session := getSession()
	defer session.Close()
	return session.DB(DB).C(c).UpdateId(id, set)
}

func UpdateById(c string, id interface{}, set B) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).UpdateId(id, set)
}

func UpdateAll(c string, query B, set B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(defaultDb).C(c).UpdateAll(query, set)
	return err
}

func DbUpdateAll(db, c string, query B, set B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(db).C(c).UpdateAll(query, set)
	return err
}

func UpsertById(c string, id string, set B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(defaultDb).C(c).UpsertId(id, set)
	return err
}

func DbUpsertById(db, c, id string, set B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(db).C(c).UpsertId(id, set)
	return err
}

func DeleteOne(c string, query B) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Remove(query)
}

func DBDeleteOne(dbName, c string, query B) error {
	session := getSession()
	defer session.Close()
	return session.DB(dbName).C(c).Remove(query)
}

func DeleteAll(c string, query B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(defaultDb).C(c).RemoveAll(query)
	return err
}

func DBDeleteAll(dbName, c string, query B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(dbName).C(c).RemoveAll(query)
	return err
}

func DeleteById(c string, id interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).RemoveId(id)
}

func GetCount(c string, query B) (int, error) {
	session := getSession()
	defer session.Close()
	return session.DB(defaultDb).C(c).Find(query).Count()
}

func GetCountByDb(dbName, c string, query B) (int, error) {
	session := getSession()
	defer session.Close()
	return session.DB(dbName).C(c).Find(query).Count()
}

func GetMaxLimitCountByDb(dbName, c string, query B, limit int) (int, error) {
	session := getSession()
	defer session.Close()
	return session.DB(dbName).C(c).Find(query).Limit(limit).Count()
}

func AggregateByDb(db string, c string, match B, group B, sort B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Pipe([]B{{"$match": match}, {"$group": group}, {"$sort": sort}}).All(result)
}

func AggregateLimitByDb(db string, c string, match B, group B, sort B, skip int, limit int, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Pipe([]B{{"$match": match}, {"$group": group}, {"$sort": sort}, {"$skip": skip}, {"$limit": limit}}).All(result)
}

func AggregateCountByDb(db string, c string, match B, group B, count string, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Pipe([]B{{"$match": match}, {"$group": group}, {"$count": count}}).One(result)
}

func DBFind(dbName, c string, query B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(dbName).C(c).Find(query).All(result)
}

func DbUpSertOne(db, c string, query B, set B) error {
	session := getSession()
	defer session.Close()
	_, err := session.DB(db).C(c).Upsert(query, set)
	if err != nil {
		return err
	}
	return nil
}

func Aggregate(db string, c string, match B, group B, project B, sort B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Pipe([]B{{"$match": match}, {"$group": group}, {"$project": project}, {"$sort": sort}}).All(result)
}

func DbFindSortFields(db string, c string, query B, field B, sort string, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Select(field).Sort(sort).All(result)
}

func DBDistinct(db string, c string, field string, query B, result interface{}) error {
	session := getSession()
	defer session.Close()
	return session.DB(db).C(c).Find(query).Distinct(field, result)
}
