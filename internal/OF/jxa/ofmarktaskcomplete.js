// Author: Mike Rhodes
// Source: https://github.com/mikerhodes/github-to-omnifocus
//
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

// Mark a task complete in OmniFocus
// Accepts a Task as JSON in an OSA_ARGS env var
// Call it:
//   set -gx OSA_ARGS '{"id": "a2g4XFUiQKm"}'
//   osascript -l JavaScript ofmarktaskcomplete.js | jq .

/**
 * @typedef {Object} OmnifocusTask
 * @property {string} id
 * @property {string} name // not used, but here to mirror type on Go side.
 */

function markTaskComplete(
    /** @type {OmnifocusTask} */ t
) {
    // @ts-ignore
    const ofApp = Application("OmniFocus")
    const task = ofApp.defaultDocument.flattenedTasks.whose({ id: t.id })[0]
    if (task) {
        // @ts-ignore
        ofApp.markComplete(task)
        return true
    }
    return false
}


ObjC.import('stdlib')
var args = JSON.parse($.getenv('OSA_ARGS'))
var out = markTaskComplete(args)
JSON.stringify(out)
