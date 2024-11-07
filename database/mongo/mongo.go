package mongo

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	InsertMany(context.Context, []interface{}) ([]interface{}, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	Aggregate(context.Context, interface{}) (Cursor, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type Cursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(context.Context, interface{}) error
}

type Client interface {
	Database(string) Database
	Connect(context.Context) error
	Disconnect(context.Context) error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	Ping(context.Context) error
}

type mongoClient struct {
	client *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	collection *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	mc *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

type nullawareDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

func (d *nullawareDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() != bsontype.Null {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if !val.CanSet() {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if err := vr.ReadNull(); err != nil {
		return err
	}

	val.Set(d.zeroValue)
	return nil
}

func GetMongoClient(connection string) (Client, error) {
	time.Local = time.UTC

	clientOptions := options.Client().ApplyURI(connection)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	return &mongoClient{client: client}, err
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.client.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) Database(dbName string) Database {
	db := mc.client.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return mc.client.UseSession(ctx, fn)
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.client.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.client.Connect(ctx)
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.client.Disconnect(ctx)
}

func (md *mongoDatabase) Collection(colName string) Collection {
	collection := md.db.Collection(colName)
	return &mongoCollection{collection: collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.db.Client()
	return &mongoClient{client: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := mc.collection.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateOne(ctx, filter, update, opts[:]...)
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.collection.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) InsertMany(ctx context.Context, documents []interface{}) ([]interface{}, error) {
	res, err := mc.collection.InsertMany(ctx, documents)
	return res.InsertedIDs, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.collection.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	findResult, err := mc.collection.Find(ctx, filter, opts...)
	return &mongoCursor{mc: findResult}, err
}

func (mc *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (Cursor, error) {
	aggregateResult, err := mc.collection.Aggregate(ctx, pipeline)

	return &mongoCursor{mc: aggregateResult}, err
}

func (mc *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateMany(ctx, filter, update, opts[:]...)
}

func (mc *mongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return mc.collection.CountDocuments(ctx, filter, opts...)
}

func (sr *mongoSingleResult) Decode(value interface{}) error {
	return sr.sr.Decode(value)
}

func (mr *mongoCursor) Close(ctx context.Context) error {
	return mr.mc.Close(ctx)
}

func (mr *mongoCursor) Next(ctx context.Context) bool {
	return mr.mc.Next(ctx)
}

func (mr *mongoCursor) Decode(v interface{}) error {
	return mr.mc.Decode(v)
}

func (mr *mongoCursor) All(ctx context.Context, result interface{}) error {
	return mr.mc.All(ctx, result)
}