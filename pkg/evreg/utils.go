// evreg is simple event registry allowing people to register for attendee
// limited slots.
// Copyright (C) 2021  Alexander Hirsch
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package evreg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func marshalToFile(filepath string, v interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	vJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	file.Write(vJSON)
	return err
}

func unmarshalFromFile(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, v)
}
