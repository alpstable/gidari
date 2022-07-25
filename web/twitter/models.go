package twitter

import "time"

// * This is a generated file, do not edit

// Annotations are details about annotations relative to the text within a Tweet.
type Annotation struct {
	// End is the end position (zero based) of the text used to annotate the Tweet. While all other end indices are
	// exclusive, this one is inclusive.
	End int `bson:"end" json:"end" sql:"end"`

	// NormalizedText is the text used to determine the annotation type.
	NormalizedText string `bson:"normalized_text" json:"normalized_text" sql:"normalized_text"`

	// Probability is the confidence score for the annotation as it correlates to the Tweet text.
	Probability float64 `bson:"probability" json:"probability" sql:"probability"`

	// Start is the start position (zero-based) of the text used to annotate the Tweet. All start indices are inclusive.
	Start int `bson:"start" json:"start" sql:"start"`

	// Type is the description of the type of entity identified when the Tweet text was interpreted.
	Type string `bson:"type" json:"type" sql:"type"`
}

// TODO
type Attachments struct {
	// MediaKeys is a list of unique identifiers of media attached to this Tweet. These identifiers use the same media key
	// format as those returned by the Media Library. You can obtain the expanded object in includes.media by adding
	// expansions=attachments.media_keys in the request's query parameter.
	MediaKeys []string `bson:"media_keys" json:"media_keys" sql:"media_keys"`

	// PollIds list of unique identifiers of polls present in the Tweets returned. These are returned as a string in order
	// to avoid complications with languages and tools that cannot handle large integers.
	PollIds []string `bson:"poll_ids" json:"poll_ids" sql:"poll_ids"`
}

// TODO
type Bookmark struct {
	// Attachments specifies the type of attachments (if any) present in this Tweet. To return this field, add
	// tweet.fields=attachments in the request's query parameter.
	Attachments *Attachments `bson:"attachments" json:"attachments" sql:"attachments"`

	// Authur is the inique identifier of this user. This is returned as a string in order to avoid complications with
	// languages and tools that cannot handle large integers. You can obtain the expanded object in includes.users by adding
	// expansions=author_id in the request's query parameter.
	AuthorID string `bson:"author_id" json:"author_id" sql:"author_id"`

	// Bookmarked indicates whether the user has removed the Bookmark of the specified Tweet. specified Tweet as a result of
	// this request. The returned value is false for a successful request. If the data has been created through a POST
	// method, Bookmarked indicates whether the user bookmarks the specified Tweet as a result of this request.
	Bookmarked bool `bson:"bookmarked" json:"bookmarked" sql:"bookmarked"`

	// ContextAnnotations are context annotations for the Tweet. To return this field, add tweet.fields=context_annotations
	// in the request's query parameter.
	ContextAnnotations []*ContextAnnotation `bson:"context_annotations" json:"context_annotations" sql:"context_annotations"`

	// ConversationID is the Tweet ID of the original Tweet of the conversation (which includes direct replies, replies of
	// replies). To return this field, add tweet.fields=conversation_id in the request's query parameter.
	ConversationID string `bson:"conversation_id" json:"conversation_id" sql:"conversation_id"`

	// CreatedAt is the creation time of the Tweet.
	CreatedAt time.Time `bson:"created_at" json:"created_at" sql:"created_at"`

	// Entities contain details about text that has a special meaning in a Tweet. To return this field, add
	// tweet.fields=entities in the request's query parameter.
	Entities *Entities `bson:"entities" json:"entities" sql:"entities"`

	// Geo contains details about the location tagged by the user in this Tweet, if they specified one. To return this
	// field, add tweet.fields=geo in the request's query parameter.
	Geo *Geo `bson:"geo" json:"geo" sql:"geo"`

	// ID is a unique identifier of this Tweet. This is returned as a string in order to avoid complications with languages
	// and tools that cannot handle large integers.
	ID string `bson:"id" json:"id" sql:"id"`

	// InReplyToUserID indicates the user ID of the parent Tweet's author. This is returned as a string in order to avoid
	// complications with languages and tools that cannot handle large integers. You can obtain the expanded object in
	// includes.users by adding expansions=in_reply_to_user_id in the request's query parameter.
	InReplyToUserID string `bson:"in_reply_to_user_id" json:"in_reply_to_user_id" sql:"in_reply_to_user_id"`

	// NonPublicMetrics are non-public engagement metrics for the Tweet at the time of the request. This is a private
	// metric, and requires the use of OAuth 2.0 User Context authentication. To return this field, add
	// tweet.fields=non_public_metrics in the request's query parameter.
	NonPublicMetrics *NonPublicMetrics `bson:"non_public_metrics" json:"non_public_metrics" sql:"non_public_metrics"`

	// OrganicMetrics are organic engagement metrics for the Tweet at the time of the request. Requires user context
	// authentication.
	OrganicMetrics *OrganicMetrics `bson:"organic_metrics" json:"organic_metrics" sql:"organic_metrics"`

	// PromotedMetrics are engagement metrics for the Tweet at the time of the request in a promoted context. Requires user
	// context authentication.
	PromotedMetrics *PromotedMetrics `bson:"promoted_metrics" json:"promoted_metrics" sql:"promoted_metrics"`

	// PublicMetrics are the engagement metrics for the Tweet at the time of the request. To return this field, add
	// tweet.fields=public_metrics in the request's query parameter.
	PublicMetrics *PublicMetrics `bson:"public_metrics" json:"public_metrics" sql:"public_metrics"`

	// ReferencedTweets is a list of Tweets this Tweet refers to. For example, if the parent Tweet is a Retweet, a Retweet
	// with comment (also known as Quoted Tweet) or a Reply, it will include the related Tweet referenced to by its parent.
	// To return this field, add tweet.fields=referenced_tweets in the request's query parameter.
	ReferencedTweets []*ReferencedTweet `bson:"referenced_tweets" json:"referenced_tweets" sql:"referenced_tweets"`

	// Text is the content of the Tweet.
	Text string `bson:"text" json:"text" sql:"text"`

	// Withheld contains withholding details for withheld content. To return this field, add tweet.fields=withheld in the
	// request's query parameter.
	Withheld          *Withheld              `bson:"withheld" json:"withheld" sql:"withheld"`
	Errors            map[string]interface{} `bson:"errors" json:"errors" sql:"errors"`
	Includes          map[string]interface{} `bson:"includes" json:"includes" sql:"includes"`
	Lang              string                 `bson:"lang" json:"lang" sql:"lang"`
	PossiblySensitive bool                   `bson:"possibly_sensitive" json:"possibly_sensitive" sql:"possibly_sensitive"`
	ReplySettings     string                 `bson:"reply_settings" json:"reply_settings" sql:"reply_settings"`
	Source            string                 `bson:"source" json:"source" sql:"source"`
}

// BookmarkWrite details the results from a bookmark write operation.
type BookmarkWrite struct {
	Data Bookmarked `bson:"data" json:"data" sql:"data"`
	Meta Meta       `bson:"meta" json:"meta" sql:"meta"`
}

// Bookmarked holds details about the status of a bookmark write.
type Bookmarked struct {
	Bookmarked bool `bson:"bookmarked" json:"Bookmarked" sql:"bookmarked"`
}

// TODO
type Bookmarks struct {
	Data []*Bookmark `bson:"data" json:"data" sql:"data"`
	Meta Meta        `bson:"meta" json:"meta" sql:"meta"`
}

// Cashtag contains details about text recognized as a Cashtag.
type Cashtag struct {
	// End is the end position (zero-based) of the recognized Cashtag within the Tweet. This end index is exclusive.
	End int `bson:"end" json:"end" sql:"end"`

	// Start is the start position (zero-based) of the recognized Cashtag within the Tweet. All start indices are inclusive.
	Start int `bson:"start" json:"start" sql:"start"`

	// Tag is the text of the Cashtag.
	Tag string `bson:"tag" json:"tag" sql:"tag"`
}

// Compliance is some recent compliance jobs.
type Compliance struct {
	// CreatedAt is the date and time when the job was created.
	CreatedAt time.Time `bson:"created_at" json:"created_at" sql:"created_at"`

	// DownloadExpiresAt the date and time until which the download URL will be available (usually 7 days from the request
	// time).
	DownloadExpiresAt time.Time `bson:"download_expires_at" json:"download_expires_at" sql:"download_expires_at"`

	// DownloadURL is the predefined location where to download the results from the compliance job. This URL is already
	// signed with an authentication key, so you will not need to pass any additional credential or header to authenticate
	// the request.
	DownloadURL string `bson:"download_url" json:"download_url" sql:"download_url"`

	// Error returns when jobs.status is failed. Specifies the reason why the job did not complete successfully.
	Error string `bson:"error" json:"error" sql:"error"`

	// ID is the unique identifier for this job.
	ID string `bson:"id" json:"id" sql:"id"`

	// Meta returns meta information about the request.
	Meta Meta `bson:"meta" json:"meta" sql:"meta"`

	// Name is the user defined job name. Only returned if specified when the job was created.
	Name string `bson:"name" json:"name" sql:"name"`

	// Status is the status of this job.
	Status Status `bson:"status" json:"status" sql:"status"`

	// Type is the type of the job, whether tweets or users.
	Type ComplianceJob `bson:"type" json:"type" sql:"type"`

	// UploadExpiresAt represents the date and time until which the upload URL will be available (usually 15 minutes from
	// the request time).
	UploadExpiresAt time.Time `bson:"upload_expires_at" json:"upload_expires_at" sql:"upload_expires_at"`

	// UploadURL is a URL representing the location where to upload IDs consumed by your app. This URL is already signed
	// with an authentication key, so you will not need to pass any additional credentials or headers to authenticate the
	// request.
	UploadURL string `bson:"upload_url" json:"upload_url" sql:"upload_url"`
}

// TODO
type ContextAnnotation struct {
	// Domain are elements which identify detailed information regarding the domain classification based on Tweet text.
	Domain Domain `bson:"domain" json:"domain" sql:"domain"`

	// Entity are elements which identify detailed information regarding the domain classification bases on Tweet text.
	Entity Entity `bson:"entity" json:"entity" sql:"entity"`
}

// TODO
type Coordinates struct {
	Coordinates []float64 `bson:"coordinates" json:"coordinates" sql:"coordinates"`
	Type        string    `bson:"type" json:"type" sql:"type"`
}

// Domain identifies detailed information regarding the domain classification based on Tweet text.
type Domain struct {
	// Description is the Long form description of domain classification.
	Description string `bson:"description" json:"description" sql:"description"`

	// ID is the numeric value of the domain.
	ID string `bson:"id" json:"id" sql:"id"`

	// Name is the domain name based on the Tweet text.
	Name string `bson:"name" json:"name" sql:"name"`
}

// Entities are details about text that has a special meaning in a Tweet. To return this field, add
// tweet.fields=entities in the request's query parameter.
type Entities struct {
	// Annotations contain details about annotations relative to the text within a Tweet.
	Annotations []*Annotation `bson:"annotations" json:"annotations" sql:"annotations"`

	// Cashtags contain details about text recognized as a Cashtag.
	Cashtags []*Cashtag `bson:"cashtags" json:"cashtags" sql:"cashtags"`

	// Hashtags contains details about text recognized as a Hashtag.
	Hashtags []*Hashtag `bson:"hashtags" json:"hashtags" sql:"hashtags"`

	// Mentions contains details about text recognized as a user mention.
	Mentions []*Mention `bson:"mentions" json:"mentions" sql:"mentions"`

	// Urls contains details about text recognized as a URL.
	Urls []*URL `bson:"urls" json:"urls" sql:"urls"`
}

// Entity identifies detailed information regarding the domain classification bases on Tweet text.
type Entity struct {
	// Description is additional information regarding referenced entity.
	Description string `bson:"description" json:"description" sql:"description"`

	// ID is a unique value which correlates to an explicitly mentioned Person, Place, Product or Organization.
	ID string `bson:"id" json:"id" sql:"id"`

	// Name is the name or reference of entity referenced in the Tweet.
	Name string `bson:"name" json:"name" sql:"name"`
}

// TODO
type Geo struct {
	Coordinates *Coordinates `bson:"coordinates" json:"coordinates" sql:"coordinates"`
	PlaceID     string       `bson:"place_id" json:"place_id" sql:"place_id"`
}

// Hashtag contains details about text recognized as a Hashtag.
type Hashtag struct {
	// End is the end position (zero-based) of the recognized Hashtag within the Tweet. This end index is exclusive.
	End int `bson:"end" json:"end" sql:"end"`

	// Start is the start position (zero-based) of the recognized Hashtag within the Tweet. All start indices are inclusive.
	Start int `bson:"start" json:"start" sql:"start"`

	// Tag is the text of the Hashtag.
	Tag string `bson:"tag" json:"tag" sql:"tag"`
}

// Mention contains details about text recognized as a user mention.
type Mention struct {
	// End is the end position (zero-based) of the recognized user mention within the Tweet. This end index is exclusive.
	End int `bson:"end" json:"end" sql:"end"`

	// Start is the start position (zero-based) of the recognized user mention within the Tweet. All start indices are
	// inclusive.
	Start int `bson:"start" json:"start" sql:"start"`

	// Username is the part of text recognized as a user mention. You can obtain the expanded object in includes.users by
	// adding expansions=entities.mentions.username in the request's query parameter.
	Username string `bson:"username" json:"username" sql:"username"`
}

// Meta holds metadata concerning requests
type Meta struct {
	ResultCount int `bson:"result_count" json:"result_count" sql:"result_count"`
}

// NonPublicMetrics are non-public engagement metrics for the Tweet at the time of the request. This is a private
// metric, and requires the use of OAuth 2.0 User Context authentication.
type NonPublicMetrics struct {
	// ImpressionCount are the number of times the Tweet has been viewed. This is a private metric, and requires the use of
	// OAuth 2.0 User Context authentication..
	ImpressionCount int `bson:"impression_count" json:"impression_count" sql:"impression_count"`

	// URLLinkClicks are the number of times a user clicks on a URL link or URL preview card in a Tweet. This is a private
	// metric, and requires the use of OAuth 2.0 User Context authentication.
	URLLinkClicks int `bson:"url_link_clicks" json:"url_link_clicks" sql:"url_link_clicks"`

	// UserProfileClicks are the number of times a user clicks the following portions of a Tweet - display name, user name,
	// profile picture. This is a private metric, and requires the use of OAuth 2.0 User Context authentication.
	UserProfileClicks int `bson:"user_profile_clicks" json:"user_profile_clicks" sql:"user_profile_clicks"`
}

// OrganicMetrics are engagement metrics for the Tweet at the time of the request. Requires user context authentication.
type OrganicMetrics struct {
	// ImpressionCount is the number of times the Tweet has been viewed organically. This is a private metric, and requires
	// the use of OAuth 2.0 User Context authentication.
	ImpressionCount int `bson:"impression_count" json:"impression_count" sql:"impression_count"`

	// LikeCount is the number of likes the Tweet has received organically.
	LikeCount int `bson:"like_count" json:"like_count" sql:"like_count"`

	// ReplyCount is the number of replies the Tweet has received organically.
	ReplyCount int `bson:"reply_count" json:"reply_count" sql:"reply_count"`

	// RetweetCountis the number of times the Tweet has been Retweeted organically.
	RetweetCount int `bson:"retweet_count" json:"retweet_count" sql:"retweet_count"`

	// URLLinkClicks is the number of times a user clicks on a URL link or URL preview card in a Tweet organically. This is
	// a private metric, and requires the use of OAuth 2.0 User Context authentication.
	URLLinkClicks int `bson:"url_link_clicks" json:"url_link_clicks" sql:"url_link_clicks"`

	// UserProfileClicks is the number of times a user clicks the following portions of a Tweet organically - display name,
	// user name, profile picture. This is a private metric, and requires the use of OAuth 2.0 User Context authentication.
	UserProfileClicks int `bson:"user_profile_clicks" json:"user_profile_clicks" sql:"user_profile_clicks"`
}

// PromotedMetrics are engagement metrics for the Tweet at the time of the request in a promoted context. Requires user
// context authentication.
type PromotedMetrics struct {
	ImpressionCount   int `bson:"impression_count" json:"impression_count" sql:"impression_count"`
	LikeCount         int `bson:"like_count" json:"like_count" sql:"like_count"`
	ReplyCount        int `bson:"reply_count" json:"reply_count" sql:"reply_count"`
	RetweetCount      int `bson:"retweet_count" json:"retweet_count" sql:"retweet_count"`
	URLLinkClicks     int `bson:"url_link_clicks" json:"url_link_clicks" sql:"url_link_clicks"`
	UserProfileClicks int `bson:"user_profile_clicks" json:"user_profile_clicks" sql:"user_profile_clicks"`
}

// PublicMetrics are the engagement metrics for the Tweet at the time of the request.
type PublicMetrics struct {
	// LikeCount is the number of Likes of this Tweet.
	LikeCount int `bson:"like_count" json:"like_count" sql:"like_count"`

	// QuoteCount is the number of times this Tweet has been Retweeted with a comment (also known as Quote).
	QuoteCount int `bson:"quote_count" json:"quote_count" sql:"quote_count"`

	// ReplyCount is the number of Replies of this Tweet.
	ReplyCount int `bson:"reply_count" json:"reply_count" sql:"reply_count"`

	// RetweetCount is the number of times this Tweet has been Retweeted.
	RetweetCount int `bson:"retweet_count" json:"retweet_count" sql:"retweet_count"`
}

// ReferencedTweet is a tweet referenced by another tweet.
type ReferencedTweet struct {
	// ID is the unique identifier of the referenced Tweet. You can obtain the expanded object in includes.tweets by adding
	// expansions=referenced_tweets.id in the request's query parameter.
	ID string `bson:"id" json:"id" sql:"id"`

	// Type indicates the type of relationship between this Tweet and the Tweet returned in the response: retweeted (this
	// Tweet is a Retweet), quoted (a Retweet with comment, also known as Quoted Tweet), or replied_to (this Tweet is a
	// reply).
	Type string `bson:"type" json:"type" sql:"type"`
}

// TODO
type Tweet struct {
	// Creation time of the Tweet. To return this field, add tweet.fields=created_at in the request's query parameter.
	CreatedAt time.Time `bson:"created_at" json:"created_at" sql:"created_at"`

	// ID is a unique identifier of this Tweet. This is returned as a string in order to avoid complications with languages
	// and tools that cannot handle large integers.
	ID string `bson:"id" json:"id" sql:"id"`

	// Text is the content of the Tweet. To return this field, add tweet.fields=text in the request's query parameter.
	Text string `bson:"text" json:"text" sql:"text"`
}

// TODO
type Tweets struct {
	Data []*Tweet `bson:"data" json:"data" sql:"data"`
	Meta Meta     `bson:"meta" json:"meta" sql:"meta"`
}

// Entity identifies detailed information regarding the domain classification bases on Tweet text.
type URL struct {
	// DisplayUrl is the URL as displayed in the Twitter client.
	DisplayURL string `bson:"display_url" json:"display_url" sql:"display_url"`

	// End is the end position (zero-based) of the recognized URL within the Tweet. This end index is exclusive.
	End int `bson:"end" json:"end" sql:"end"`

	// ExpandedUrl is the The fully resolved URL.
	ExpandedURL string `bson:"expanded_url" json:"expanded_url" sql:"expanded_url"`

	// Start is the start position (zero-based) of the recognized URL within the Tweet. All start indices are inclusive.
	Start int `bson:"start" json:"start" sql:"start"`

	// UnwoundUrl is the full destination URL.
	UnwoundURL string `bson:"unwound_url" json:"unwound_url" sql:"unwound_url"`

	// Url is the URL in the format tweeted by the user.
	URL string `bson:"url" json:"url" sql:"url"`
}

// TODO
type User struct {
	// ID is a unique identifier of this user. This is returned as a string in order to avoid complications with languages
	// and tools that cannot handle large integers.
	ID string `bson:"id" json:"id" sql:"id"`
}

// Withheld contains withholding details for withheld content.
type Withheld struct {
	// Copyright indicates if the content is being withheld for on the basis of copyright infringement.
	Copyright bool `bson:"copyright" json:"copyright" sql:"copyright"`

	// CountryCodes is a list of countries where this content is not available.
	CountryCodes []string `bson:"country_codes" json:"country_codes" sql:"country_codes"`

	// Scope indicates whether the content being withheld is a Tweet or a user.
	Scope WithholdingScope `bson:"scope" json:"scope" sql:"scope"`
}
