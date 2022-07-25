package twitter

import "strings"

// * This is a generated file, do not edit

type ComplianceJob string

const (
	ComplianceJobTweets ComplianceJob = "tweets"
	ComplianceJobUsers  ComplianceJob = "users"
)

// String will convert a ComplianceJob into a string.
func (ComplianceJob *ComplianceJob) String() string {
	if ComplianceJob != nil {
		return string(*ComplianceJob)
	}
	return ""
}

type Expansion string
type Expansions []Expansion

const (
	ExpansionAttachmentsPollIds         Expansion = "attachments.poll_ids"
	ExpansionAttachmentsMediaKeys       Expansion = "attachments.media_keys"
	ExpansionAuthorID                   Expansion = "author_id"
	ExpansionEntitiesMentionsUsername   Expansion = "entities.mentions.username"
	ExpansionGeoPlaceID                 Expansion = "geo.place_id"
	ExpansionInReplyToUserID            Expansion = "in_reply_to_user_id"
	ExpansionReferencedTweetsID         Expansion = "referenced_tweets.id"
	ExpansionReferencedTweetsIDAuthorID Expansion = "referenced_tweets.id.author_id"
)

// String will convert a Expansion into a string.
func (Expansion *Expansion) String() string {
	if Expansion != nil {
		return string(*Expansion)
	}
	return ""
}

// String will convert a slice of Expansion into a CSV.
func (Expansions *Expansions) String() string {
	var str string
	if Expansions != nil {
		slice := []string{}
		for _, val := range *Expansions {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

type MediaField string
type MediaFields []MediaField

const (
	MediaFieldDurationMs       MediaField = "duration_ms"
	MediaFieldHeight           MediaField = "height"
	MediaFieldMediaKey         MediaField = "media_key"
	MediaFieldPreviewImageURL  MediaField = "preview_image_url"
	MediaFieldType             MediaField = "type"
	MediaFieldURL              MediaField = "url"
	MediaFieldWidth            MediaField = "width"
	MediaFieldPublicMetrics    MediaField = "public_metrics"
	MediaFieldNonPublicMetrics MediaField = "non_public_metrics"
	MediaFieldOrganicMetrics   MediaField = "organic_metrics"
	MediaFieldPromotedMetrics  MediaField = "promoted_metrics"
	MediaFieldAltText          MediaField = "alt_text"
)

// String will convert a MediaField into a string.
func (MediaField *MediaField) String() string {
	if MediaField != nil {
		return string(*MediaField)
	}
	return ""
}

// String will convert a slice of MediaField into a CSV.
func (MediaFields *MediaFields) String() string {
	var str string
	if MediaFields != nil {
		slice := []string{}
		for _, val := range *MediaFields {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

type PlaceField string
type PlaceFields []PlaceField

const (
	PlaceFieldContainedWithin PlaceField = "contained_within"
	PlaceFieldCountry         PlaceField = "country"
	PlaceFieldCountryCode     PlaceField = "country_code"
	PlaceFieldFullName        PlaceField = "full_name"
	PlaceFieldGeo             PlaceField = "geo"
	PlaceFieldID              PlaceField = "id"
	PlaceFieldName            PlaceField = "name"
	PlaceFieldPlaceType       PlaceField = "place_type"
)

// String will convert a PlaceField into a string.
func (PlaceField *PlaceField) String() string {
	if PlaceField != nil {
		return string(*PlaceField)
	}
	return ""
}

// String will convert a slice of PlaceField into a CSV.
func (PlaceFields *PlaceFields) String() string {
	var str string
	if PlaceFields != nil {
		slice := []string{}
		for _, val := range *PlaceFields {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

type PollField string
type PollFields []PollField

const (
	PollFieldDurationMinutes PollField = "duration_minutes"
	PollFieldEndDatetime     PollField = "end_datetime"
	PollFieldID              PollField = "id"
	PollFieldOptions         PollField = "options"
	PollFieldVotingStatus    PollField = "voting_status"
)

// String will convert a PollField into a string.
func (PollField *PollField) String() string {
	if PollField != nil {
		return string(*PollField)
	}
	return ""
}

// String will convert a slice of PollField into a CSV.
func (PollFields *PollFields) String() string {
	var str string
	if PollFields != nil {
		slice := []string{}
		for _, val := range *PollFields {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

type ReferencedTweetType string
type ReferencedTweetTypes []ReferencedTweetType

const (
	ReferencedTweetTypeRetweeted ReferencedTweetType = "retweeted"
	ReferencedTweetTypeQuoted    ReferencedTweetType = "quoted"
	ReferencedTweetTypeRepliedTo ReferencedTweetType = "replied_to"
)

// String will convert a ReferencedTweetType into a string.
func (ReferencedTweetType *ReferencedTweetType) String() string {
	if ReferencedTweetType != nil {
		return string(*ReferencedTweetType)
	}
	return ""
}

// String will convert a slice of ReferencedTweetType into a CSV.
func (ReferencedTweetTypes *ReferencedTweetTypes) String() string {
	var str string
	if ReferencedTweetTypes != nil {
		slice := []string{}
		for _, val := range *ReferencedTweetTypes {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

type Status string

const (
	StatusOnline     Status = "online"
	StatusOffline    Status = "offline"
	StatusInternal   Status = "internal"
	StatusDelisted   Status = "delisted"
	StatusPending    Status = "pending"
	StatusCreating   Status = "creating"
	StatusReady      Status = "ready"
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusFailed     Status = "failed"
	StatusComplete   Status = "complete"
)

// String will convert a Status into a string.
func (Status *Status) String() string {
	if Status != nil {
		return string(*Status)
	}
	return ""
}

type TweetField string
type TweetFields []TweetField

const (
	TweetFieldAttachments        TweetField = "attachments"
	TweetFieldAuthorID           TweetField = "author_id"
	TweetFieldContextAnnotations TweetField = "context_annotations"
	TweetFieldConversationID     TweetField = "conversation_id"
	TweetFieldCreatedAt          TweetField = "created_at"
	TweetFieldEntities           TweetField = "entities"
	TweetFieldGeo                TweetField = "geo"
	TweetFieldID                 TweetField = "id"
	TweetFieldInReplyToUserID    TweetField = "in_reply_to_user_id"
	TweetFieldLang               TweetField = "lang"
	TweetFieldNonPublicMetrics   TweetField = "non_public_metrics"
	TweetFieldPublicMetrics      TweetField = "public_metrics"
	TweetFieldOrganicMetrics     TweetField = "organic_metrics"
	TweetFieldPromotedMetrics    TweetField = "promoted_metrics"
	TweetFieldPossiblySensitive  TweetField = "possibly_sensitive"
	TweetFieldReferencedTweets   TweetField = "referenced_tweets"
	TweetFieldReplySettings      TweetField = "reply_settings"
	TweetFieldSource             TweetField = "source"
	TweetFieldText               TweetField = "text"
	TweetFieldWithheld           TweetField = "withheld"
)

// String will convert a TweetField into a string.
func (TweetField *TweetField) String() string {
	if TweetField != nil {
		return string(*TweetField)
	}
	return ""
}

// String will convert a slice of TweetField into a CSV.
func (TweetFields *TweetFields) String() string {
	var str string
	if TweetFields != nil {
		slice := []string{}
		for _, val := range *TweetFields {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

type UserField string
type UserFields []UserField

const (
	UserFieldCreatedAt       UserField = "created_at"
	UserFieldDescription     UserField = "description"
	UserFieldEntities        UserField = "entities"
	UserFieldID              UserField = "id"
	UserFieldLocation        UserField = "location"
	UserFieldName            UserField = "name"
	UserFieldPinnedTweetID   UserField = "pinned_tweet_id"
	UserFieldProfileImageURL UserField = "profile_image_url"
	UserFieldProtected       UserField = "protected"
	UserFieldPublicMetrics   UserField = "public_metrics"
	UserFieldURL             UserField = "url"
	UserFieldUsername        UserField = "username"
	UserFieldVerified        UserField = "verified"
	UserFieldWithheld        UserField = "withheld"
)

// String will convert a UserField into a string.
func (UserField *UserField) String() string {
	if UserField != nil {
		return string(*UserField)
	}
	return ""
}

// String will convert a slice of UserField into a CSV.
func (UserFields *UserFields) String() string {
	var str string
	if UserFields != nil {
		slice := []string{}
		for _, val := range *UserFields {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}

// WithholdingScope indicates whether the content being withheld is a Tweet or a user.
type WithholdingScope string
type WithholdingScopes []WithholdingScope

const (
	WithholdingScopeTweet WithholdingScope = "tweet"
	WithholdingScopeUser  WithholdingScope = "user"
)

// String will convert a WithholdingScope into a string.
func (WithholdingScope *WithholdingScope) String() string {
	if WithholdingScope != nil {
		return string(*WithholdingScope)
	}
	return ""
}

// String will convert a slice of WithholdingScope into a CSV.
func (WithholdingScopes *WithholdingScopes) String() string {
	var str string
	if WithholdingScopes != nil {
		slice := []string{}
		for _, val := range *WithholdingScopes {
			slice = append(slice, val.String())
		}
		str = strings.Join(slice, ",")
	}
	return str
}
