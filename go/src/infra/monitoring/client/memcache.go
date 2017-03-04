package client

import (
	"encoding/json"
	"fmt"
	"infra/monitoring/messages"

	"github.com/luci/gae/service/memcache"
	"github.com/luci/luci-go/common/logging"

	"golang.org/x/net/context"
)

// NewMemcacheReader returns a reader which caches builds in memcache.
// Any calls using this reader should have memcache installed in the context.
func NewMemcacheReader(r readerType) readerType {
	return &memcacheReader{r: r}
}

type memcacheReader struct {
	r readerType
}

type res struct {
	b   *messages.Build
	err error
}

func (m *memcacheReader) Build(ctx context.Context, master *messages.MasterLocation, builder string, buildNum int64) (*messages.Build, error) {
	key := fmt.Sprintf("%s/%s/%d", master.Name(), builder, buildNum)
	itm := memcache.NewItem(ctx, key)
	err := memcache.Get(ctx, itm)
	if err == memcache.ErrCacheMiss {
		res, err := m.r.Build(ctx, master, builder, buildNum)
		if err != nil {
			return nil, err
		}

		if res.Finished {
			// FIXME: Maybe don't use json for serialization format.
			data, err := json.Marshal(res)
			if err != nil {
				logging.Warningf(ctx, "failed to serialize build data, when saving to memcache: %s", err)
			}
			itm.SetValue(data)
			if err = memcache.Set(ctx, itm); err != nil {
				logging.Warningf(ctx, "failed to save build data to memcache: %s", err)
			}
		} else {
			logging.Debugf(ctx, "not caching %s in memcache, as it's still pending", key)
		}
		return res, nil
	}

	dec := &messages.Build{}
	if err = json.Unmarshal(itm.Value(), dec); err != nil {
		return nil, fmt.Errorf("failed to decode data in memcache (data probably corrupt). key %s err %s", key, err)
	}

	return dec, nil
}

func (m *memcacheReader) LatestBuilds(ctx context.Context, master *messages.MasterLocation, build string) ([]*messages.Build, error) {
	return m.r.LatestBuilds(ctx, master, build)
}

func (m *memcacheReader) TestResults(ctx context.Context, masterName *messages.MasterLocation, builderName, stepName string, buildNumber int64) (*messages.TestResults, error) {
	return m.r.TestResults(ctx, masterName, builderName, stepName, buildNumber)
}

func (m *memcacheReader) BuildExtract(ctx context.Context, master *messages.MasterLocation) (*messages.BuildExtract, error) {
	return m.r.BuildExtract(ctx, master)
}

func (m *memcacheReader) StdioForStep(ctx context.Context, master *messages.MasterLocation, builder, step string, buildNum int64) ([]string, error) {
	return m.r.StdioForStep(ctx, master, builder, step, buildNum)
}

func (m *memcacheReader) CrbugItems(ctx context.Context, label string) ([]messages.CrbugItem, error) {
	return m.r.CrbugItems(ctx, label)
}

func (m *memcacheReader) Findit(ctx context.Context, master *messages.MasterLocation, builder string, buildNum int64, failedSteps []string) ([]*messages.FinditResult, error) {
	return m.r.Findit(ctx, master, builder, buildNum, failedSteps)
}
