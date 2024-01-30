package repository

import (
	"article-service/logger"
	"article-service/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func collArticles() *mongo.Collection {
	return client.Database("article-alpha").Collection("article")
}

// GetArticles retrieves a list of articles from the MongoDB collection based on the provided filter.
// If no filter is provided, it retrieves all articles from the collection.
// It returns the list of articles or an error.
func GetArticles(ctx context.Context, filter bson.D) ([]*model.Article, error) {
	var articles []*model.Article

	// Set an empty filter as default if none is provided
	if filter == nil {
		filter = bson.D{}
	}

	// Query the MongoDB collection with the provided filter
	cursor, err := collArticles().Find(ctx, filter)
	if err != nil {
		logger.Logger.Debugf("GetArticles - Error %s", err)
		return nil, err
	}

	// Decode the query result into a list of articles
	if err := cursor.All(ctx, &articles); err != nil {
		logger.Logger.Debugf("GetArticles - Error %s", err)
		return nil, err
	}

	logger.Logger.Debugf("GetArticles - Found %d articles", len(articles))
	return articles, nil
}

// AddArticle adds a new article to the MongoDB collection.
// It assigns a new ObjectID to the article before inserting it into the collection.
// It returns an error if any.
func AddArticle(ctx context.Context, article *model.Article) error {
	// Generate a new ObjectID for the article
	article.ID = primitive.NewObjectID().Hex()

	// Insert the article into the MongoDB collection
	_, err := collArticles().InsertOne(ctx, article)
	if err != nil {
		logger.Logger.Debugf("AddArticle - Error %s", err)
		return err
	}

	logger.Logger.Debugf("AddArticle - Article inserted in the database with ID %s", article.ID)
	return nil
}

// DeleteArticle deletes an article from the MongoDB collection based on its ID.
// It returns an error if any.
func DeleteArticle(ctx context.Context, articleID string) error {
	var err error

	// Convert the article ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(articleID)
	if err != nil {
		logger.Logger.Debugf("DeleteArticle - Error %s", err)
		return err
	}

	// Delete the article from the MongoDB collection using its ObjectID
	_, err = collArticles().DeleteOne(ctx, bson.D{{"_id", objectID.Hex()}})
	if err != nil {
		logger.Logger.Debugf("DeleteArticle - Error %s", err)
		return err
	}

	logger.Logger.Debugf("DeleteArticle - Article with ID %s deleted", articleID)
	return nil
}
