package relations

import "go.mongodb.org/mongo-driver/bson/primitive"

type dtoRelations struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserId      string             `bson:"user_id"`
	FriendsIds  []string           `bson:"friends"`
	BlockedIds  []string           `bson:"blocked"`
	IncomingIds []string           `bson:"incoming"`
	OutgoingIds []string           `bson:"outgoing"`
}
