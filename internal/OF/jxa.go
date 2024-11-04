// Author: Mike Rhodes
// Modified by: Trevor Piltch
// Source: https://github.com/mikerhodes/github-to-omnifocus

/*
 * Copyright 2020 Mike Rhodes, https://dx13.co.uk/
Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
*/

package omnifocus

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

// ItemsForQuery returns a list of items from Omnifocus that
// match the passed query.
func ItemsForQuery(q ItemQuery) ([]Item, error) {
	jsCode, _ := jxa.ReadFile("jxa/oftasksforprojectwithtag.js")
	args, _ := json.Marshal(q)

	out, err := executeScript(jsCode, args)
	if err != nil {
		return []Item{}, err
	}

	items := []Item{}
	err = json.Unmarshal(out, &items)
	if err != nil {
		return []Item{}, err
	}

	return items, nil
}

// MarkOmniFocusItemComplete marks a Item as complete. It only requires the
// id field to be set.
func MarkOmnifocusItemComplete(i Item) error {
	jsCode, _ := jxa.ReadFile("jxa/ofmarktaskcomplete.js")
	args, _ := json.Marshal(i)

	_, err := executeScript(jsCode, args)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// EnsureTagExists creates a tag in OmniFocus if it doesn't already exist.
func EnsureTagExists(tag Tag) error {
	jsCode, _ := jxa.ReadFile("jxa/ofensuretagexists.js")
	args, _ := json.Marshal(tag)

	_, err := executeScript(jsCode, args)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// AddNewOmnifocusItem adds a new OmniFocus Item
func AddNewOmnifocusItem(t NewOmniFocusItem) (Item, error) {
	jsCode, _ := jxa.ReadFile("jxa/ofaddnewtask.js")
	args, _ := json.Marshal(t)

	out, err := executeScript(jsCode, args)
	if err != nil {
		return Item{}, err
	}

	item := Item{}
	err = json.Unmarshal(out, &item)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

// executeScript runs jsCode passing it args as input, and returns the
// output of the command.
func executeScript(jsCode []byte, args []byte) ([]byte, error) {
	// All scripts expect a JSON object passed in via the
	// OSA_ARGS environment variable. The script itself is
	// passed into osascript via stdin. The script outputs
	// a JSON document over stdout.

	cmd := exec.Command("/usr/bin/osascript", "-l", "JavaScript")

	cmd.Env = append(os.Environ(),
		"OSA_ARGS="+string(args),
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, string(jsCode))
		if err != nil {
			// should never fail
			log.Fatal(err)
		}
	}()

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
