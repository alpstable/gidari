package repository

import (
	"context"
	"time"

	"github.com/alpine-hodler/driver/data/proto"
	"github.com/alpine-hodler/driver/internal/query"
	"github.com/alpine-hodler/driver/tools"
	"github.com/alpine-hodler/driver/web/polygon"
	"google.golang.org/protobuf/types/known/structpb"
)

// Polygon is the repository wrapper for creating CRUD operations on a storage object.
type Polygon interface {
	tools.GenericStorage

	Name() string
	ReadBarMinutes(context.Context, string, bool, time.Time, time.Time, *proto.ReadResponse) error
	UpsertBarMinutes(context.Context, []*polygon.BarResult, *proto.CreateResponse) ([]*polygon.BarResult, error)
}

type p struct{ *storage }

// NewPolygon will return a new service for performing CRUD operations on the Coinbase Pro candles web API.
func NewPolygon(_ context.Context, r tools.GenericStorage) Polygon {
	stg := new(storage)
	stg.r = newStorage(r)
	return &p{storage: stg}
}

func (svc *p) Name() string {
	return "Polygon"
}

// ReadBarMinutes will query storage for Polygon Bars to the minute-granularity for a given ticker, adjusted value,
// start time, and end time.
func (svc *p) ReadBarMinutes(ctx context.Context, ticker string, adjusted bool, start time.Time, end time.Time,
	rsp *proto.ReadResponse) error {

	req := new(proto.ReadRequest)
	req.ReaderBuilder = query.PolygonBarMinutesReadBuilder.Bytes()
	req.Table = "bar_minutes"
	if err := tools.AssignReadRequired(req, "ticker", ticker); err != nil {
		return err
	}
	if err := tools.AssignReadRequired(req, "adjusted", adjusted); err != nil {
		return err
	}
	if err := tools.AssignReadRequired(req, "start", start.Unix()); err != nil {
		return err
	}
	if err := tools.AssignReadRequired(req, "end", end.Unix()); err != nil {
		return err
	}
	return svc.r.Read(ctx, req, rsp)
}

// UpsertBarMinutes will upsert bar results upto a minutes to the "bar_minutes" resource.
func (svc *p) UpsertBarMinutes(ctx context.Context, bars []*polygon.BarResult,
	rsp *proto.CreateResponse) ([]*polygon.BarResult, error) {

	var records []*structpb.Struct
	if err := tools.MakeRecordsRequest(bars, &records); err != nil {
		return nil, err
	}
	req := new(proto.UpsertRequest)
	req.Table = "bar_minutes"
	req.Records = records
	return bars, svc.r.Upsert(ctx, req, rsp)
}
