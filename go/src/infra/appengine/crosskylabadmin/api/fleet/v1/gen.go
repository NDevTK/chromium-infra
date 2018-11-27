// Copyright 2018 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate cproto

//go:generate mockgen -source tracker.pb.go -destination tracker.mock.pb.go -package fleet
//go:generate mockgen -source tasker.pb.go -destination tasker.mock.pb.go -package fleet
//go:generate mockgen -source inventory.pb.go -destination inventory.mock.pb.go -package fleet

//go:generate svcdec -type TrackerServer
//go:generate svcdec -type TaskerServer
//go:generate svcdec -type InventoryServer

// Package fleet contains service definitions for fleet management in
// crosskylabadmin.
package fleet
