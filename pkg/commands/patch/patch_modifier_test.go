package patch

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const simpleDiff = `diff --git a/filename b/filename
index dcd3485..1ba5540 100644
--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-orange
+grape
 ...
 ...
 ...
`

const addNewlineToEndOfFile = `diff --git a/filename b/filename
index 80a73f1..e48a11c 100644
--- a/filename
+++ b/filename
@@ -60,4 +60,4 @@ grape
 ...
 ...
 ...
-last line
\ No newline at end of file
+last line
`

const removeNewlinefromEndOfFile = `diff --git a/filename b/filename
index e48a11c..80a73f1 100644
--- a/filename
+++ b/filename
@@ -60,4 +60,4 @@ grape
 ...
 ...
 ...
-last line
+last line
\ No newline at end of file
`

const twoHunks = `diff --git a/filename b/filename
index e48a11c..b2ab81b 100644
--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-grape
+orange
 ...
 ...
 ...
@@ -8,6 +8,8 @@ grape
 ...
 ...
 ...
+pear
+lemon
 ...
 ...
 ...
`

const twoChangesInOneHunk = `diff --git a/filename b/filename
index 9320895..6d79956 100644
--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-grape
+kiwi
 orange
-pear
+banana
 lemon
`

const newFile = `diff --git a/newfile b/newfile
new file mode 100644
index 0000000..4e680cc
--- /dev/null
+++ b/newfile
@@ -0,0 +1,3 @@
+apple
+orange
+grape
`

const addNewlineToPreviouslyEmptyFile = `diff --git a/newfile b/newfile
index e69de29..c6568ea 100644
--- a/newfile
+++ b/newfile
@@ -0,0 +1 @@
+new line
\ No newline at end of file
`

const exampleHunk = `@@ -1,5 +1,5 @@
 apple
-grape
+orange
...
...
...
`

// TestModifyPatchForRange is a function.
func TestModifyPatchForRange(t *testing.T) {
	type scenario struct {
		testName       string
		filename       string
		diffText       string
		firstLineIndex int
		lastLineIndex  int
		reverse        bool
		expected       string
	}

	scenarios := []scenario{
		{
			testName:       "nothing selected",
			filename:       "filename",
			firstLineIndex: -1,
			lastLineIndex:  -1,
			diffText:       simpleDiff,
			expected:       "",
		},
		{
			testName:       "only context selected",
			filename:       "filename",
			firstLineIndex: 5,
			lastLineIndex:  5,
			diffText:       simpleDiff,
			expected:       "",
		},
		{
			testName:       "whole range selected",
			filename:       "filename",
			firstLineIndex: 0,
			lastLineIndex:  11,
			diffText:       simpleDiff,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-orange
+grape
 ...
 ...
 ...
`,
		},
		{
			testName:       "only removal selected",
			filename:       "filename",
			firstLineIndex: 6,
			lastLineIndex:  6,
			diffText:       simpleDiff,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,4 @@
 apple
-orange
 ...
 ...
 ...
`,
		},
		{
			testName:       "only addition selected",
			filename:       "filename",
			firstLineIndex: 7,
			lastLineIndex:  7,
			diffText:       simpleDiff,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,6 @@
 apple
 orange
+grape
 ...
 ...
 ...
`,
		},
		{
			testName:       "range that extends beyond diff bounds",
			filename:       "filename",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       simpleDiff,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-orange
+grape
 ...
 ...
 ...
`,
		},
		{
			testName:       "add newline to end of file",
			filename:       "filename",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       addNewlineToEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,4 +60,4 @@ grape
 ...
 ...
 ...
-last line
\ No newline at end of file
+last line
`,
		},
		{
			testName:       "add newline to end of file, reversed",
			filename:       "filename",
			firstLineIndex: -100,
			lastLineIndex:  100,
			reverse:        true,
			diffText:       addNewlineToEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,4 +60,4 @@ grape
 ...
 ...
 ...
-last line
\ No newline at end of file
+last line
`,
		},
		{
			testName:       "remove newline from end of file",
			filename:       "filename",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       removeNewlinefromEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,4 +60,4 @@ grape
 ...
 ...
 ...
-last line
+last line
\ No newline at end of file
`,
		},
		{
			testName:       "remove newline from end of file, reversed",
			filename:       "filename",
			firstLineIndex: -100,
			lastLineIndex:  100,
			reverse:        true,
			diffText:       removeNewlinefromEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,4 +60,4 @@ grape
 ...
 ...
 ...
-last line
+last line
\ No newline at end of file
`,
		},
		{
			testName:       "remove newline from end of file, removal only",
			filename:       "filename",
			firstLineIndex: 8,
			lastLineIndex:  8,
			diffText:       removeNewlinefromEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,4 +60,3 @@ grape
 ...
 ...
 ...
-last line
`,
		},
		{
			testName:       "remove newline from end of file, removal only, reversed",
			filename:       "filename",
			firstLineIndex: 8,
			lastLineIndex:  8,
			reverse:        true,
			diffText:       removeNewlinefromEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,5 +60,4 @@ grape
 ...
 ...
 ...
-last line
 last line
\ No newline at end of file
`,
		},
		{
			testName:       "remove newline from end of file, addition only",
			filename:       "filename",
			firstLineIndex: 9,
			lastLineIndex:  9,
			diffText:       removeNewlinefromEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,4 +60,5 @@ grape
 ...
 ...
 ...
 last line
+last line
\ No newline at end of file
`,
		},
		{
			testName:       "remove newline from end of file, addition only, reversed",
			filename:       "filename",
			firstLineIndex: 9,
			lastLineIndex:  9,
			reverse:        true,
			diffText:       removeNewlinefromEndOfFile,
			expected: `--- a/filename
+++ b/filename
@@ -60,3 +60,4 @@ grape
 ...
 ...
 ...
+last line
\ No newline at end of file
`,
		},
		{
			testName:       "staging two whole hunks",
			filename:       "filename",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       twoHunks,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-grape
+orange
 ...
 ...
 ...
@@ -8,6 +8,8 @@ grape
 ...
 ...
 ...
+pear
+lemon
 ...
 ...
 ...
`,
		},
		{
			testName:       "staging part of both hunks",
			filename:       "filename",
			firstLineIndex: 7,
			lastLineIndex:  15,
			diffText:       twoHunks,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,6 @@
 apple
 grape
+orange
 ...
 ...
 ...
@@ -8,6 +9,7 @@ grape
 ...
 ...
 ...
+pear
 ...
 ...
 ...
`,
		},
		{
			testName:       "adding a new file",
			filename:       "newfile",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       newFile,
			expected: `--- a/newfile
+++ b/newfile
@@ -0,0 +1,3 @@
+apple
+orange
+grape
`,
		},
		{
			testName:       "adding part of a new file",
			filename:       "newfile",
			firstLineIndex: 6,
			lastLineIndex:  7,
			diffText:       newFile,
			expected: `--- a/newfile
+++ b/newfile
@@ -0,0 +1,2 @@
+apple
+orange
`,
		},
		{
			testName:       "adding a new line to a previously empty file",
			filename:       "newfile",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       addNewlineToPreviouslyEmptyFile,
			expected: `--- a/newfile
+++ b/newfile
@@ -0,0 +1,1 @@
+new line
\ No newline at end of file
`,
		},
		{
			testName:       "adding a new line to a previously empty file, reversed",
			filename:       "newfile",
			firstLineIndex: -100,
			lastLineIndex:  100,
			diffText:       addNewlineToPreviouslyEmptyFile,
			reverse:        true,
			expected: `--- a/newfile
+++ b/newfile
@@ -0,0 +1,1 @@
+new line
\ No newline at end of file
`,
		},
		{
			testName:       "adding part of a hunk",
			filename:       "filename",
			firstLineIndex: 6,
			lastLineIndex:  7,
			reverse:        false,
			diffText:       twoChangesInOneHunk,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-grape
+kiwi
 orange
 pear
 lemon
`,
		},
		{
			testName:       "adding part of a hunk, reverse",
			filename:       "filename",
			firstLineIndex: 6,
			lastLineIndex:  7,
			reverse:        true,
			diffText:       twoChangesInOneHunk,
			expected: `--- a/filename
+++ b/filename
@@ -1,5 +1,5 @@
 apple
-grape
+kiwi
 orange
 banana
 lemon
`,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			result := ModifiedPatchForRange(nil, s.filename, s.diffText, s.firstLineIndex, s.lastLineIndex,
				PatchOptions{
					Reverse:            s.reverse,
					KeepOriginalHeader: false,
				})
			if !assert.Equal(t, s.expected, result) {
				fmt.Println(result)
			}
		})
	}
}

func TestLineNumberOfLine(t *testing.T) {
	type scenario struct {
		testName string
		hunk     *PatchHunk
		idx      int
		expected int
	}

	scenarios := []scenario{
		{
			testName: "nothing selected",
			hunk:     newHunk(strings.SplitAfter(exampleHunk, "\n"), 10),
			idx:      15,
			expected: 3,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			result := s.hunk.LineNumberOfLine(s.idx)
			if !assert.Equal(t, s.expected, result) {
				fmt.Println(result)
			}
		})
	}
}
