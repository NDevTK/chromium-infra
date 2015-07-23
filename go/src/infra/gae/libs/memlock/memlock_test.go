// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package memlock

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/luci/gae/filter/featureBreaker"
	"github.com/luci/gae/impl/memory"
	"github.com/luci/gae/service/memcache"
	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/clock/testclock"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func init() {
	delay = time.Millisecond
	memcacheLockTime = time.Millisecond * 4
}

type getBlockerFilter struct {
	memcache.Interface
	sync.Mutex
}

func (f *getBlockerFilter) Get(key string) (memcache.Item, error) {
	f.Lock()
	defer f.Unlock()
	return f.Interface.Get(key)
}

func TestSimple(t *testing.T) {
	// TODO(riannucci): Mock time.After so that we don't have to delay for real.

	const key = memlockKeyPrefix + "testkey"

	Convey("basic locking", t, func() {
		start := time.Date(1986, time.October, 26, 1, 20, 00, 00, time.UTC)
		ctx, clk := testclock.UseTime(context.Background(), start)
		blocker := make(chan struct{})
		clk.SetTimerCallback(func(clock.Timer) {
			clk.Add(delay)
			select {
			case blocker <- struct{}{}:
			default:
			}
		})
		waitFalse := func(ctx context.Context) {
		loop:
			for {
				select {
				case <-blocker:
					continue
				case <-ctx.Done():
					break loop
				}
			}
		}

		ctx, fb := featureBreaker.FilterMC(memory.Use(ctx), nil)
		mc := memcache.Get(ctx)

		Convey("fails to acquire when memcache is down", func() {
			fb.BreakFeatures(nil, "Add")
			err := TryWithLock(ctx, "testkey", "id", func(context.Context) error {
				// should never reach here
				So(false, ShouldBeTrue)
				return nil
			})
			So(err, ShouldEqual, ErrFailedToLock)
		})

		Convey("returns the inner error", func() {
			toRet := fmt.Errorf("sup")
			err := TryWithLock(ctx, "testkey", "id", func(context.Context) error {
				return toRet
			})
			So(err, ShouldEqual, toRet)
		})

		Convey("returns the error", func() {
			toRet := fmt.Errorf("sup")
			err := TryWithLock(ctx, "testkey", "id", func(context.Context) error {
				return toRet
			})
			So(err, ShouldEqual, toRet)
		})

		Convey("can acquire when empty", func() {
			err := TryWithLock(ctx, "testkey", "id", func(ctx context.Context) error {
				isDone := func() bool {
					select {
					case <-ctx.Done():
						return true
					default:
						return false
					}
				}

				So(isDone(), ShouldBeFalse)

				Convey("waiting for a while keeps refreshing the lock", func() {
					// simulate waiting for 64*delay time, and ensuring that checkLoop
					// runs that many times.
					for i := 0; i < 64; i++ {
						<-blocker
						clk.Add(delay)
					}
					So(isDone(), ShouldBeFalse)
				})

				Convey("but sometimes we might lose it", func() {
					Convey("because it was evicted", func() {
						mc.Delete(key)
						clk.Add(memcacheLockTime)
						waitFalse(ctx)
					})

					Convey("or because of service issues", func() {
						fb.BreakFeatures(nil, "CompareAndSwap")
						waitFalse(ctx)
					})
				})
				return nil
			})
			So(err, ShouldBeNil)
		})

		Convey("can lose it when it gets stolen", func() {
			gbf := &getBlockerFilter{}
			ctx = memcache.AddFilters(ctx, func(_ context.Context, mc memcache.Interface) memcache.Interface {
				gbf.Interface = mc
				return gbf
			})
			mc = memcache.Get(ctx)
			err := TryWithLock(ctx, "testkey", "id", func(ctx context.Context) error {
				// simulate waiting for 64*delay time, and ensuring that checkLoop
				// runs that many times.
				for i := 0; i < 64; i++ {
					<-blocker
					clk.Add(delay)
				}
				gbf.Lock()
				mc.Set(mc.NewItem(key).SetValue([]byte("wat")))
				gbf.Unlock()
				waitFalse(ctx)
				return nil
			})
			So(err, ShouldBeNil)
		})

		Convey("an empty context id is an error", func() {
			So(TryWithLock(ctx, "testkey", "", nil), ShouldEqual, ErrEmptyClientID)
		})
	})
}
