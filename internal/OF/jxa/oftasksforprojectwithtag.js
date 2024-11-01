// Author: Mike Rhodes
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

//
// Return the tasks for a project having a given tag
// Accepts a TaskQuery as JSON in an OSA_ARGS env var.
// Call it:
//   set -gx OSA_ARGS '{"projectName": "GitHub Notifications", "tags": ["github"]}'
//   osascript -l JavaScript oftasksforprojectwithtag.js | jq .
// Returns JSON array:
// [
//     {
//       "id": "iAKv1Uo8XqW",
//       "name": "cloudant/techspec-documents#257 Document modernize search project progress"
//     }, ...
// ]

/**
 * @typedef {Object} TaskQuery
 * @property {string} projectName
 * @property {string[]} tags
 */

function tasksForProjectWithTag(
    /** @type {TaskQuery} */ query
) {
    // @ts-ignore
    const ofApp = Application("OmniFocus")
    const ofDoc = ofApp.defaultDocument
    const project = ofDoc.flattenedProjects
        .whose({ name: query.projectName })[0];

    const tagFoundOrCreated = charTag => {
        const
            tags = ofDoc.flattenedTags.whose({
                name: charTag
            }),
            oTag = ofApp.Tag({
                name: charTag
            });
        return tags.length === 0 ? (
            (
                ofDoc.tags.push(oTag),
                oTag
            )
        ) : tags()[0]
    }

    const ofTags = query.tags.map((t) => {
        return tagFoundOrCreated(t)
    })

    // const ofAppTag = tagFoundOrCreated(appTag)
    // const ofTypeTag = tagFoundOrCreated(typeTag)

    return project.tasks()
        .filter((task) => task.completed() === false)
        .filter((task) => {
            // Task must have all tags
            const tags = task.tags()
            for (var i = 0; i < ofTags.length; i++) {
                if (!tags.some(tag => tag.id() == ofTags[i].id())) {
                    return false
                }
            }
            return true
        })
        .map((task) => {
            return { "id": task.id(), "name": task.name() };
        });
}

ObjC.import('stdlib')
var args = JSON.parse($.getenv('OSA_ARGS'))
var out = tasksForProjectWithTag(args)
JSON.stringify(out)
