package repository

import (
	"context"
	"errors"
	"time"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "chat_direct"
)

type MongoChatDirect struct {
	Id        bson.ObjectID              `bson:"_id"`
	SenderId  bson.ObjectID              `bson:"sender_id"`
	User1Id   bson.ObjectID              `bson:"user1Id"`
	User2Id   bson.ObjectID              `bson:"user2Id"`
	Content   string                     `bson:"content"`
	Reactions []*MongoChatDirectReaction `bson:"reactions"`
	CreatedAt time.Time                  `bson:"created_at"`
	UpdatedAt time.Time                  `bson:"updated_at"`
}

type MongoChatDirectReaction struct {
	Id       bson.ObjectID   `bson:"_id"`
	Users    []bson.ObjectID `bson:"user_ids"`
	Reaction string          `bson:"reaction"`
}

type MongoChatDirectRepositoryDeps struct {
	uberdig.In
	Mapper mapper.Mapper[*MongoChatDirect, *entity.ChatDirect]
	Client *mongo.Client
}

type MongoChatDirectRepository struct {
	deps MongoChatDirectRepositoryDeps
}

func NewMongoChatDirectRepository(deps MongoChatDirectRepositoryDeps) ChatDirectRepository {
	return &MongoChatDirectRepository{deps: deps}
}

func (m MongoChatDirectRepository) Create(
	ctx context.Context,
	chatDirect *entity.ChatDirect,
) error {
	if chatDirect.Id == "" {
		chatDirect.Id = entity.ChatDirectId(bson.NewObjectID().Hex())
	}
	now := time.Now()
	if chatDirect.CreatedAt.IsZero() {
		chatDirect.CreatedAt = now
		chatDirect.UpdatedAt = now
	}
	if chatDirect.UpdatedAt.IsZero() {
		chatDirect.UpdatedAt = now
	}

	mongoChatDirect, err := m.deps.Mapper.MapFromEntity(chatDirect)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoChatDirect)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChatDirectRepository) GetById(
	ctx context.Context,
	chatDirectId entity.ChatDirectId,
) (*entity.ChatDirect, error) {
	objectId, err := bson.ObjectIDFromHex(string(chatDirectId))
	if err != nil {
		return nil, err
	}

	var mongoChatDirect MongoChatDirect
	err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": objectId}).Decode(&mongoChatDirect)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, nil // Pas de message trouvé
		}
		return nil, err
	}

	return m.deps.Mapper.MapToEntity(&mongoChatDirect)
}

func (m MongoChatDirectRepository) ListRecentChats(
	ctx context.Context,
	userId user_entity.UserId,
) ([]*entity.ChatDirect, error) {
	userObjectId, err := bson.ObjectIDFromHex(string(userId))
	if err != nil {
		return nil, err
	}

	pipeline := mongo2.Pipeline{
		// Étape 0 : Filtrer les conversations où l'utilisateur est impliqué
		{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "user1Id", Value: userObjectId}},
				bson.D{{Key: "user2Id", Value: userObjectId}},
			}},
		}}},
		// Étape 1 : Ajouter un champ `sortedIds` qui normalise l'ordre des IDs
		{{Key: "$addFields", Value: bson.D{
			{Key: "sortedIds", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{{Key: "$lt", Value: bson.A{"$user1Id", "$user2Id"}}}},
					{Key: "then", Value: bson.A{"$user1Id", "$user2Id"}},
					{Key: "else", Value: bson.A{"$user2Id", "$user1Id"}},
				}},
			}},
		}}},
		// Étape 2 : Grouper par `sortedIds` en prenant le message le plus récent
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$sortedIds"},
			{Key: "latestMessage", Value: bson.D{{Key: "$last", Value: "$$ROOT"}}},
		}}},
		// Étape 3 : Trier par `updated_at` du dernier message
		{{Key: "$sort", Value: bson.D{{Key: "latestMessage.updated_at", Value: -1}}}},
	}

	cursor, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var groupedMongoChatDirects []struct {
		LatestMessage MongoChatDirect `bson:"latestMessage"`
	}

	if err = cursor.All(ctx, &groupedMongoChatDirects); err != nil {
		return nil, err
	}

	chatDirects := make([]*entity.ChatDirect, len(groupedMongoChatDirects))
	for i, groupedMongoChatDirect := range groupedMongoChatDirects {
		chatDirect, err := m.deps.Mapper.MapToEntity(&groupedMongoChatDirect.LatestMessage)
		if err != nil {
			return nil, err
		}
		chatDirects[i] = chatDirect
	}

	return chatDirects, nil
}

func (m MongoChatDirectRepository) IsFirstMessage(
	ctx context.Context,
	user1Id, user2Id user_entity.UserId,
) (bool, error) {
	user1IdHex, err := bson.ObjectIDFromHex(string(user1Id))
	if err != nil {
		return false, err
	}

	user2IdHex, err := bson.ObjectIDFromHex(string(user2Id))
	if err != nil {
		return false, err
	}

	filter := bson.D{
		{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "user1Id", Value: user1IdHex}, {Key: "user2Id", Value: user2IdHex}},
				bson.D{{Key: "user1Id", Value: user2IdHex}, {Key: "user2Id", Value: user1IdHex}},
			},
		},
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (m MongoChatDirectRepository) GetLastMessage(
	ctx context.Context,
	user1Id, user2Id user_entity.UserId,
) (*entity.ChatDirect, error) {
	user1IdHex, err := bson.ObjectIDFromHex(string(user1Id))
	if err != nil {
		return nil, err
	}

	user2IdHex, err := bson.ObjectIDFromHex(string(user2Id))
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "user1Id", Value: user1IdHex}, {Key: "user2Id", Value: user2IdHex}},
				bson.D{{Key: "user1Id", Value: user2IdHex}, {Key: "user2Id", Value: user1IdHex}},
			},
		},
	}

	opts := options.FindOne().
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // Dernier message

	var mongoChatDirect MongoChatDirect
	err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, filter, opts).Decode(&mongoChatDirect)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, nil // Pas de messages trouvés
		}
		return nil, err
	}

	return m.deps.Mapper.MapToEntity(&mongoChatDirect)
}

//nolint:revive
func (m MongoChatDirectRepository) ListByUser(
	ctx context.Context,
	user1Id, user2Id user_entity.UserId,
	params ListByUserQueryParams,
) ([]*entity.ChatDirect, error) {
	user1IdHex, err := bson.ObjectIDFromHex(string(user1Id))
	if err != nil {
		return nil, err
	}

	user2IdHex, err := bson.ObjectIDFromHex(string(user2Id))
	if err != nil {
		return nil, err
	}

	// Filter must match both ways
	filter := bson.D{
		{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "user1Id", Value: user1IdHex}, {Key: "user2Id", Value: user2IdHex}},
				bson.D{{Key: "user1Id", Value: user2IdHex}, {Key: "user2Id", Value: user1IdHex}},
			},
		},
	}

	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	var messages []*entity.ChatDirect
	var aroundMessage MongoChatDirect

	// Si AroundMessageId est défini, on récupère le message cible
	if params.AroundMessageId != "" {
		aroundObjectId, err := bson.ObjectIDFromHex(string(params.AroundMessageId))
		if err != nil {
			return nil, err
		}

		aroundMessageFilter := append(filter, bson.E{Key: "_id", Value: aroundObjectId})

		err = collection.FindOne(ctx, aroundMessageFilter).Decode(&aroundMessage)
		if err != nil {
			return nil, err
		}

		// Récupération des messages avant le message cible
		beforeFilter := append(
			filter,
			bson.E{Key: "created_at", Value: bson.M{"$lt": aroundMessage.CreatedAt}},
		)
		beforeOpts := options.Find().
			SetSort(bson.D{{Key: "created_at", Value: -1}}).
			// Ordre décroissant pour prendre les plus récents en premier
			SetLimit(int64(params.Limit / 2)) // On récupère la moitié avant

		beforeCursor, err := collection.Find(ctx, beforeFilter, beforeOpts)
		if err != nil {
			return nil, err
		}
		defer beforeCursor.Close(ctx)

		var beforeMessages []*entity.ChatDirect
		for beforeCursor.Next(ctx) {
			var mongoMessage MongoChatDirect
			if err := beforeCursor.Decode(&mongoMessage); err != nil {
				return nil, err
			}
			message, err := m.deps.Mapper.MapToEntity(&mongoMessage)
			if err != nil {
				return nil, err
			}
			beforeMessages = append(beforeMessages, message)
		}

		// Récupération des messages après le message cible
		afterFilter := append(
			filter,
			bson.E{Key: "created_at", Value: bson.M{"$gt": aroundMessage.CreatedAt}},
		)
		afterOpts := options.Find().
			SetSort(bson.D{{Key: "created_at", Value: 1}}).
			// Ordre croissant pour prendre les plus anciens en premier
			SetLimit(int64(params.Limit / 2)) // On récupère la moitié après

		afterCursor, err := collection.Find(ctx, afterFilter, afterOpts)
		if err != nil {
			return nil, err
		}
		defer afterCursor.Close(ctx)

		var afterMessages []*entity.ChatDirect
		for afterCursor.Next(ctx) {
			var mongoMessage MongoChatDirect
			if err := afterCursor.Decode(&mongoMessage); err != nil {
				return nil, err
			}
			message, err := m.deps.Mapper.MapToEntity(&mongoMessage)
			if err != nil {
				return nil, err
			}
			afterMessages = append(afterMessages, message)
		}

		// On assemble le tout dans l'ordre chronologique
		for i := len(beforeMessages) - 1; i >= 0; i-- { // On remet l'ordre chronologique
			messages = append(messages, beforeMessages[i])
		}
		channelMessage, err := m.deps.Mapper.MapToEntity(&aroundMessage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, channelMessage) // Ajout du message cible
		messages = append(messages, afterMessages...)

		return messages, nil
	}

	// Sinon, on applique les filtres classiques
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // Trier par date décroissante
		SetLimit(int64(params.Limit))

	if !params.Before.Equal(time.Time{}) {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$lt": params.Before}})
	} else if !params.After.Equal(time.Time{}) {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$gt": params.After}})
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoChatDirects []*MongoChatDirect
	if err = cursor.All(ctx, &mongoChatDirects); err != nil {
		return nil, err
	}

	chatDirects := make([]*entity.ChatDirect, len(mongoChatDirects))
	for i, mongoChatDirect := range mongoChatDirects {
		chatDirect, err := m.deps.Mapper.MapToEntity(mongoChatDirect)
		if err != nil {
			return nil, err
		}
		chatDirects[i] = chatDirect
	}

	return chatDirects, nil
}

// ToggleReaction toggles the reaction of a user to a message. (If the user has already reacted, it will remove the reaction, otherwise it will add the reaction.)
//
//nolint:revive
func (m MongoChatDirectRepository) ToggleReaction(
	ctx context.Context,
	messageId entity.ChatDirectId,
	userId user_entity.UserId,
	reaction string,
) (added bool, err error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	messageObjectId, err := bson.ObjectIDFromHex(string(messageId))
	if err != nil {
		return false, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return false, err
	}

	var message MongoChatDirect
	err = collection.FindOne(ctx, bson.M{"_id": messageObjectId}).Decode(&message)
	if err != nil {
		return false, err
	}

	updatedReactions := make([]*MongoChatDirectReaction, 0)
	found := false
	removed := false

	for _, r := range message.Reactions {
		if r.Reaction == reaction {
			found = true
			updatedUsers := make([]bson.ObjectID, 0, len(r.Users))
			for _, uid := range r.Users {
				if uid.Hex() == userId.String() {
					removed = true
					continue
				}
				updatedUsers = append(updatedUsers, uid)
			}
			if len(updatedUsers) > 0 {
				r.Users = updatedUsers
				updatedReactions = append(updatedReactions, r)
			}
		} else {
			updatedReactions = append(updatedReactions, r)
		}
	}

	if !found {
		// Add new reaction if not found
		updatedReactions = append(updatedReactions, &MongoChatDirectReaction{
			Id:       bson.NewObjectID(),
			Users:    []bson.ObjectID{userObjectId},
			Reaction: reaction,
		})
	} else if !removed {
		// Add user to existing reaction
		for _, r := range updatedReactions {
			if r.Reaction == reaction {
				r.Users = append(r.Users, userObjectId)
				break
			}
		}
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": messageObjectId},
		bson.M{"$set": bson.M{"reactions": updatedReactions}},
	)
	if err != nil {
		return false, err
	}

	return !removed, nil
}

func (m *MongoChatDirectRepository) ListAllMessagesByUser(ctx context.Context, userId user_entity.UserId) ([]*entity.ChatDirect, error) {
	userIdHex, err := bson.ObjectIDFromHex(string(userId))
	if err != nil {
		return nil, err
	}

	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	filter := bson.M{
		"$or": []bson.M{
			{"user1Id": userIdHex},
			{"user2Id": userIdHex},
		},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoChatDirects []*MongoChatDirect
	if err := cursor.All(ctx, &mongoChatDirects); err != nil {
		return nil, err
	}

	chatDirects := make([]*entity.ChatDirect, len(mongoChatDirects))
	for i, mongoChatDirect := range mongoChatDirects {
		chatDirect, err := m.deps.Mapper.MapToEntity(mongoChatDirect)
		if err != nil {
			return nil, err
		}
		chatDirects[i] = chatDirect
	}
	return chatDirects, nil
}

func (m *MongoChatDirectRepository) DeleteMessage(ctx context.Context, chatDirectId entity.ChatDirectId) error {
	chatDirectObjectId, err := bson.ObjectIDFromHex(string(chatDirectId))
	if err != nil {
		return err
	}

	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": chatDirectObjectId})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoChatDirectRepository) UpdateMessage(ctx context.Context, msg *entity.ChatDirect) error {
	if msg.User2Id.IsAfter(msg.User1Id) {
		msg.User1Id, msg.User2Id = msg.User2Id, msg.User1Id
	}

	mongoChatDirect, err := m.deps.Mapper.MapFromEntity(msg)
	if err != nil {
		return err
	}

	mongoChatDirect.UpdatedAt = time.Now()

	chatDirectObjectId, err := bson.ObjectIDFromHex(string(msg.Id))
	if err != nil {
		return err
	}

	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": chatDirectObjectId},
		bson.M{"$set": mongoChatDirect},
	)
	if err != nil {
		return err
	}

	return nil
}
