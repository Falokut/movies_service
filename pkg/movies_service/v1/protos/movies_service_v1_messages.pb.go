// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: movies_service_v1_messages.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetMovieRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MovieID int32 `protobuf:"varint,1,opt,name=movieID,json=movie_id,proto3" json:"movieID,omitempty"`
}

func (x *GetMovieRequest) Reset() {
	*x = GetMovieRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMovieRequest) ProtoMessage() {}

func (x *GetMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMovieRequest.ProtoReflect.Descriptor instead.
func (*GetMovieRequest) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{0}
}

func (x *GetMovieRequest) GetMovieID() int32 {
	if x != nil {
		return x.MovieID
	}
	return 0
}

// for multiple values use ',' separator
type GetMoviesPreviewRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MoviesIDs    *string `protobuf:"bytes,1,opt,name=moviesIDs,json=movies_ids,proto3,oneof" json:"moviesIDs,omitempty"`
	GenresIDs    *string `protobuf:"bytes,2,opt,name=genresIDs,json=genres_ids,proto3,oneof" json:"genresIDs,omitempty"`
	CountriesIDs *string `protobuf:"bytes,3,opt,name=countriesIDs,json=country_ids,proto3,oneof" json:"countriesIDs,omitempty"`
	Title        *string `protobuf:"bytes,4,opt,name=title,proto3,oneof" json:"title,omitempty"`
	AgeRatings   *string `protobuf:"bytes,5,opt,name=ageRatings,json=age_ratings,proto3,oneof" json:"ageRatings,omitempty"`
	// if limit = 0, will be used default limit = 10, if bigger than 100, will be
	// used max limit = 100
	Limit  uint32 `protobuf:"varint,6,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset uint32 `protobuf:"varint,7,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetMoviesPreviewRequest) Reset() {
	*x = GetMoviesPreviewRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMoviesPreviewRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMoviesPreviewRequest) ProtoMessage() {}

func (x *GetMoviesPreviewRequest) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMoviesPreviewRequest.ProtoReflect.Descriptor instead.
func (*GetMoviesPreviewRequest) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{1}
}

func (x *GetMoviesPreviewRequest) GetMoviesIDs() string {
	if x != nil && x.MoviesIDs != nil {
		return *x.MoviesIDs
	}
	return ""
}

func (x *GetMoviesPreviewRequest) GetGenresIDs() string {
	if x != nil && x.GenresIDs != nil {
		return *x.GenresIDs
	}
	return ""
}

func (x *GetMoviesPreviewRequest) GetCountriesIDs() string {
	if x != nil && x.CountriesIDs != nil {
		return *x.CountriesIDs
	}
	return ""
}

func (x *GetMoviesPreviewRequest) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *GetMoviesPreviewRequest) GetAgeRatings() string {
	if x != nil && x.AgeRatings != nil {
		return *x.AgeRatings
	}
	return ""
}

func (x *GetMoviesPreviewRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetMoviesPreviewRequest) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type MoviePreview struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortDescription string `protobuf:"bytes,1,opt,name=shortDescription,json=short_description,proto3" json:"shortDescription,omitempty"`
	TitleRU          string `protobuf:"bytes,2,opt,name=titleRU,json=title_ru,proto3" json:"titleRU,omitempty"`
	TitleEN          string `protobuf:"bytes,3,opt,name=titleEN,json=title_en,proto3" json:"titleEN,omitempty"`
	// movie duration in minutes
	Duration         int32    `protobuf:"varint,4,opt,name=duration,proto3" json:"duration,omitempty"`
	PreviewPosterURL string   `protobuf:"bytes,5,opt,name=previewPosterURL,json=preview_poster_url,proto3" json:"previewPosterURL,omitempty"`
	Countries        []string `protobuf:"bytes,6,rep,name=countries,proto3" json:"countries,omitempty"`
	Genres           []string `protobuf:"bytes,7,rep,name=genres,proto3" json:"genres,omitempty"`
	ReleaseYear      int32    `protobuf:"varint,8,opt,name=releaseYear,json=release_year,proto3" json:"releaseYear,omitempty"`
	AgeRating        string   `protobuf:"bytes,9,opt,name=ageRating,json=age_rating,proto3" json:"ageRating,omitempty"`
}

func (x *MoviePreview) Reset() {
	*x = MoviePreview{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoviePreview) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoviePreview) ProtoMessage() {}

func (x *MoviePreview) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoviePreview.ProtoReflect.Descriptor instead.
func (*MoviePreview) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{2}
}

func (x *MoviePreview) GetShortDescription() string {
	if x != nil {
		return x.ShortDescription
	}
	return ""
}

func (x *MoviePreview) GetTitleRU() string {
	if x != nil {
		return x.TitleRU
	}
	return ""
}

func (x *MoviePreview) GetTitleEN() string {
	if x != nil {
		return x.TitleEN
	}
	return ""
}

func (x *MoviePreview) GetDuration() int32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *MoviePreview) GetPreviewPosterURL() string {
	if x != nil {
		return x.PreviewPosterURL
	}
	return ""
}

func (x *MoviePreview) GetCountries() []string {
	if x != nil {
		return x.Countries
	}
	return nil
}

func (x *MoviePreview) GetGenres() []string {
	if x != nil {
		return x.Genres
	}
	return nil
}

func (x *MoviePreview) GetReleaseYear() int32 {
	if x != nil {
		return x.ReleaseYear
	}
	return 0
}

func (x *MoviePreview) GetAgeRating() string {
	if x != nil {
		return x.AgeRating
	}
	return ""
}

type MoviesPreview struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Movies map[int32]*MoviePreview `protobuf:"bytes,1,rep,name=movies,proto3" json:"movies,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MoviesPreview) Reset() {
	*x = MoviesPreview{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoviesPreview) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoviesPreview) ProtoMessage() {}

func (x *MoviesPreview) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoviesPreview.ProtoReflect.Descriptor instead.
func (*MoviesPreview) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{3}
}

func (x *MoviesPreview) GetMovies() map[int32]*MoviePreview {
	if x != nil {
		return x.Movies
	}
	return nil
}

type AgeRatings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ratings []string `protobuf:"bytes,1,rep,name=ratings,proto3" json:"ratings,omitempty"`
}

func (x *AgeRatings) Reset() {
	*x = AgeRatings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgeRatings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgeRatings) ProtoMessage() {}

func (x *AgeRatings) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgeRatings.ProtoReflect.Descriptor instead.
func (*AgeRatings) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{4}
}

func (x *AgeRatings) GetRatings() []string {
	if x != nil {
		return x.Ratings
	}
	return nil
}

type Country struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Country) Reset() {
	*x = Country{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Country) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Country) ProtoMessage() {}

func (x *Country) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Country.ProtoReflect.Descriptor instead.
func (*Country) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{5}
}

func (x *Country) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Country) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Countries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Countries []*Country `protobuf:"bytes,1,rep,name=countries,proto3" json:"countries,omitempty"`
}

func (x *Countries) Reset() {
	*x = Countries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Countries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Countries) ProtoMessage() {}

func (x *Countries) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Countries.ProtoReflect.Descriptor instead.
func (*Countries) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{6}
}

func (x *Countries) GetCountries() []*Country {
	if x != nil {
		return x.Countries
	}
	return nil
}

type Genre struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Genre) Reset() {
	*x = Genre{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Genre) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Genre) ProtoMessage() {}

func (x *Genre) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Genre.ProtoReflect.Descriptor instead.
func (*Genre) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{7}
}

func (x *Genre) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Genre) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Genres struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Genres []*Genre `protobuf:"bytes,1,rep,name=genres,proto3" json:"genres,omitempty"`
}

func (x *Genres) Reset() {
	*x = Genres{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Genres) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Genres) ProtoMessage() {}

func (x *Genres) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Genres.ProtoReflect.Descriptor instead.
func (*Genres) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{8}
}

func (x *Genres) GetGenres() []*Genre {
	if x != nil {
		return x.Genres
	}
	return nil
}

type Movie struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string   `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	TitleRU     string   `protobuf:"bytes,2,opt,name=titleRU,json=title_ru,proto3" json:"titleRU,omitempty"`
	TitleEN     string   `protobuf:"bytes,3,opt,name=titleEN,json=title_en,proto3" json:"titleEN,omitempty"`
	Genres      []string `protobuf:"bytes,4,rep,name=genres,json=genres_names,proto3" json:"genres,omitempty"`
	// movie duration in minutes
	Duration      int32    `protobuf:"varint,5,opt,name=duration,proto3" json:"duration,omitempty"`
	Countries     []string `protobuf:"bytes,6,rep,name=countries,json=countres_names,proto3" json:"countries,omitempty"`
	PosterURL     string   `protobuf:"bytes,7,opt,name=posterURL,json=poster_url,proto3" json:"posterURL,omitempty"`
	BackgroundURL string   `protobuf:"bytes,8,opt,name=backgroundURL,json=background_url,proto3" json:"backgroundURL,omitempty"`
	ReleaseYear   int32    `protobuf:"varint,9,opt,name=releaseYear,json=release_year,proto3" json:"releaseYear,omitempty"`
	AgeRating     string   `protobuf:"bytes,10,opt,name=ageRating,json=age_rating,proto3" json:"ageRating,omitempty"`
}

func (x *Movie) Reset() {
	*x = Movie{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Movie) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Movie) ProtoMessage() {}

func (x *Movie) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Movie.ProtoReflect.Descriptor instead.
func (*Movie) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{9}
}

func (x *Movie) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Movie) GetTitleRU() string {
	if x != nil {
		return x.TitleRU
	}
	return ""
}

func (x *Movie) GetTitleEN() string {
	if x != nil {
		return x.TitleEN
	}
	return ""
}

func (x *Movie) GetGenres() []string {
	if x != nil {
		return x.Genres
	}
	return nil
}

func (x *Movie) GetDuration() int32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *Movie) GetCountries() []string {
	if x != nil {
		return x.Countries
	}
	return nil
}

func (x *Movie) GetPosterURL() string {
	if x != nil {
		return x.PosterURL
	}
	return ""
}

func (x *Movie) GetBackgroundURL() string {
	if x != nil {
		return x.BackgroundURL
	}
	return ""
}

func (x *Movie) GetReleaseYear() int32 {
	if x != nil {
		return x.ReleaseYear
	}
	return 0
}

func (x *Movie) GetAgeRating() string {
	if x != nil {
		return x.AgeRating
	}
	return ""
}

type GetMoviesPreviewByIDsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// for multiple values use ',' separator
	MoviesIDs string `protobuf:"bytes,1,opt,name=moviesIDs,json=movies_ids,proto3" json:"moviesIDs,omitempty"`
}

func (x *GetMoviesPreviewByIDsRequest) Reset() {
	*x = GetMoviesPreviewByIDsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMoviesPreviewByIDsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMoviesPreviewByIDsRequest) ProtoMessage() {}

func (x *GetMoviesPreviewByIDsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMoviesPreviewByIDsRequest.ProtoReflect.Descriptor instead.
func (*GetMoviesPreviewByIDsRequest) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{10}
}

func (x *GetMoviesPreviewByIDsRequest) GetMoviesIDs() string {
	if x != nil {
		return x.MoviesIDs
	}
	return ""
}

type UserErrorMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UserErrorMessage) Reset() {
	*x = UserErrorMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movies_service_v1_messages_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserErrorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserErrorMessage) ProtoMessage() {}

func (x *UserErrorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_movies_service_v1_messages_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserErrorMessage.ProtoReflect.Descriptor instead.
func (*UserErrorMessage) Descriptor() ([]byte, []int) {
	return file_movies_service_v1_messages_proto_rawDescGZIP(), []int{11}
}

func (x *UserErrorMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_movies_service_v1_messages_proto protoreflect.FileDescriptor

var file_movies_service_v1_messages_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x76, 0x31, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x22, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x69, 0x64,
	0x22, 0xbe, 0x02, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x50, 0x72,
	0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x09,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x49, 0x44, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x0a, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x69, 0x64, 0x73, 0x88, 0x01, 0x01,
	0x12, 0x22, 0x0a, 0x09, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x49, 0x44, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0a, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x64,
	0x73, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x49, 0x44, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0b, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x73, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x0a, 0x61, 0x67, 0x65, 0x52, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x0b, 0x61,
	0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x49, 0x44, 0x73, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x67, 0x65,
	0x6e, 0x72, 0x65, 0x73, 0x49, 0x44, 0x73, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x49, 0x44, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x73, 0x22, 0xb3, 0x02, 0x0a, 0x0c, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x50, 0x72, 0x65, 0x76, 0x69,
	0x65, 0x77, 0x12, 0x2b, 0x0a, 0x10, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x19, 0x0a, 0x07, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x52, 0x55, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x5f, 0x72, 0x75, 0x12, 0x19, 0x0a, 0x07, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x45, 0x4e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x5f, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x2c, 0x0a, 0x10, 0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x50, 0x6f, 0x73, 0x74,
	0x65, 0x72, 0x55, 0x52, 0x4c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x65,
	0x76, 0x69, 0x65, 0x77, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x12,
	0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x67,
	0x65, 0x6e, 0x72, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0b, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x59, 0x65, 0x61, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x65,
	0x61, 0x73, 0x65, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x12, 0x1d, 0x0a, 0x09, 0x61, 0x67, 0x65, 0x52,
	0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x67, 0x65,
	0x5f, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x22, 0xab, 0x01, 0x0a, 0x0d, 0x4d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x50, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x41, 0x0a, 0x06, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x50, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x1a, 0x57, 0x0a, 0x0b,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x32, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x6f,
	0x76, 0x69, 0x65, 0x50, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x26, 0x0a, 0x0a, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x2d, 0x0a,
	0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x42, 0x0a, 0x09,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x09, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x22, 0x2b, 0x0a, 0x05, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x37, 0x0a,
	0x06, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x72, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x52, 0x06,
	0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x22, 0xc4, 0x02, 0x0a, 0x05, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x07, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x52, 0x55, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x5f, 0x72, 0x75, 0x12, 0x19, 0x0a,
	0x07, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x45, 0x4e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x5f, 0x65, 0x6e, 0x12, 0x1c, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x72,
	0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x65, 0x73, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x55,
	0x52, 0x4c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x75, 0x72, 0x6c, 0x12, 0x25, 0x0a, 0x0d, 0x62, 0x61, 0x63, 0x6b, 0x67, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x55, 0x52, 0x4c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x62, 0x61, 0x63,
	0x6b, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x75, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0b, 0x72,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x59, 0x65, 0x61, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x12, 0x1d,
	0x0a, 0x09, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x3d, 0x0a,
	0x1c, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x50, 0x72, 0x65, 0x76, 0x69, 0x65,
	0x77, 0x42, 0x79, 0x49, 0x44, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x09, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x49, 0x44, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x5f, 0x69, 0x64, 0x73, 0x22, 0x2c, 0x0a, 0x10,
	0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x1a, 0x5a, 0x18, 0x6d, 0x6f,
	0x76, 0x69, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_movies_service_v1_messages_proto_rawDescOnce sync.Once
	file_movies_service_v1_messages_proto_rawDescData = file_movies_service_v1_messages_proto_rawDesc
)

func file_movies_service_v1_messages_proto_rawDescGZIP() []byte {
	file_movies_service_v1_messages_proto_rawDescOnce.Do(func() {
		file_movies_service_v1_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_movies_service_v1_messages_proto_rawDescData)
	})
	return file_movies_service_v1_messages_proto_rawDescData
}

var file_movies_service_v1_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_movies_service_v1_messages_proto_goTypes = []interface{}{
	(*GetMovieRequest)(nil),              // 0: movies_service.GetMovieRequest
	(*GetMoviesPreviewRequest)(nil),      // 1: movies_service.GetMoviesPreviewRequest
	(*MoviePreview)(nil),                 // 2: movies_service.MoviePreview
	(*MoviesPreview)(nil),                // 3: movies_service.MoviesPreview
	(*AgeRatings)(nil),                   // 4: movies_service.AgeRatings
	(*Country)(nil),                      // 5: movies_service.Country
	(*Countries)(nil),                    // 6: movies_service.Countries
	(*Genre)(nil),                        // 7: movies_service.Genre
	(*Genres)(nil),                       // 8: movies_service.Genres
	(*Movie)(nil),                        // 9: movies_service.Movie
	(*GetMoviesPreviewByIDsRequest)(nil), // 10: movies_service.GetMoviesPreviewByIDsRequest
	(*UserErrorMessage)(nil),             // 11: movies_service.UserErrorMessage
	nil,                                  // 12: movies_service.MoviesPreview.MoviesEntry
}
var file_movies_service_v1_messages_proto_depIdxs = []int32{
	12, // 0: movies_service.MoviesPreview.movies:type_name -> movies_service.MoviesPreview.MoviesEntry
	5,  // 1: movies_service.Countries.countries:type_name -> movies_service.Country
	7,  // 2: movies_service.Genres.genres:type_name -> movies_service.Genre
	2,  // 3: movies_service.MoviesPreview.MoviesEntry.value:type_name -> movies_service.MoviePreview
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_movies_service_v1_messages_proto_init() }
func file_movies_service_v1_messages_proto_init() {
	if File_movies_service_v1_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_movies_service_v1_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMovieRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMoviesPreviewRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MoviePreview); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MoviesPreview); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgeRatings); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Country); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Countries); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Genre); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Genres); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Movie); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMoviesPreviewByIDsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_movies_service_v1_messages_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserErrorMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_movies_service_v1_messages_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_movies_service_v1_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_movies_service_v1_messages_proto_goTypes,
		DependencyIndexes: file_movies_service_v1_messages_proto_depIdxs,
		MessageInfos:      file_movies_service_v1_messages_proto_msgTypes,
	}.Build()
	File_movies_service_v1_messages_proto = out.File
	file_movies_service_v1_messages_proto_rawDesc = nil
	file_movies_service_v1_messages_proto_goTypes = nil
	file_movies_service_v1_messages_proto_depIdxs = nil
}
