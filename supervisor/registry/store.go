//  Copyright (c) 2018 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"fmt"
	"github.com/mosuka/blast/supervisor/store"
)

type StoreConstructor func(config map[string]interface{}) (store.Store, error)

type StoreRegistry map[string]StoreConstructor

var (
	Stores = make(StoreRegistry, 0)
)

func RegisterStore(name string, constructor StoreConstructor) {
	_, exists := Stores[name]
	if exists {
		panic(fmt.Errorf("attempted to register duplicate store named '%s'", name))
	}
	Stores[name] = constructor
}

func StoreConstructorByName(name string) StoreConstructor {
	return Stores[name]
}
