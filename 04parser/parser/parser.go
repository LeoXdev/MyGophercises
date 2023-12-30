// Package parser exposes a single function to parse all HTML links
// present on a HTML file (by path).
//
// - HTML DOM parsing via linear approach.
// - Still misses nested links.
//
// The responsibility of writing the output goes to the caller client.
// Tests (parser_test.go) will be implemented as part of the bonus requirements.
package parser

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// link is a struct holding the required fields for the exercise.
/* newLink := link{
	href: "/roughly",
	text: "like this input",
} */
type link struct {
	Href string
	Text string
}
// links holds multiple links objects.
type links []link

// ParseLinks searchs for links inside <a> tags on a given HTML file and returns them.
func ParseLinks(path string) (links, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("an error ocurred while opening the file: " + err.Error())
	}
	tokenizer := html.NewTokenizer(file)

	result := links{}
	var newLink link
	var inProcessLinks stack

	// writing tells if a link object is being written.
	/*
		changes -> true, after an opening tag <a> 's been found
		changes -> false, after a closing tag <a> 's been found
	*/
	var writing bool = false

	// Begin HTML tokenization
	for {
		tokenType := tokenizer.Next()
		// EOF
		if tokenType == html.ErrorToken {
			break
			//return nil, errors.New("an error ocurred while tokenizing the file")
		}
		if tokenType == html.CommentToken {
			continue
		}

		// Process the current token.
		text := tokenizer.Text()
		tag, hasAtts := tokenizer.TagName()
		iattKey, iattVal, _ := tokenizer.TagAttr()

		// Remove all ocurrences of newline chars
		text = removeByte(text, 10)

		// When finding an opening <a> tag
		// *<a> tags without a href won't be considerated
		if strings.EqualFold(string(tag), "a") && !writing && hasAtts {
			// Keep iterating until the current attribute key is "href", this'll ensure
			// that iattVal holds the url
			// There may still be an opening <a> tag with other attributes other than href,
			// to cover this, we'll loop around 170 times (the total number of attributes in HTMl);
			// Source: https://www.howtocodeschool.com/2019/10/list-of-all-html-attributes-and-their-function.html
			writing = true

			ctr := 0
			for {
				// If there's not a href attribute on the <a> tag,
				if string(iattKey) == "href" || ctr == 170 {
					break
				}
				iattKey, iattVal, _ = tokenizer.TagAttr()
				ctr++
			}
			// If the ctr reached the 170 value means that the href wasn't found
			if ctr != 170 {
				newLink.Href = string(iattVal)
			}
			//newLink.href = ""
			inProcessLinks.Push(newLink)
		} else if strings.EqualFold(string(tag), "a") && writing {
			// When finding a closing <a> tag
			writing = false
			res, err := inProcessLinks.Pop()
			if err != nil {
				return nil, errors.New("accesing links stack")
			}
			result = append(result, res)
			// Reset the newLink values after appending the object to the list
			newLink = link{}
		} else if writing {
			// When finding a token and we want to write to the newLink.text
			writtingTo, err := inProcessLinks.Peek()
			if err != nil {
				return nil, errors.New("accesing links stack")
			}
			// Avoid appending blank spaces
			if text != nil {
				writtingTo.Text += string(text)
			}
			
			// Overwrite the top element in the stack (after writing to the text attribute)
			inProcessLinks.Pop()
			inProcessLinks.Push(writtingTo)
		}
		// Concatenate each text value found until another a is found
		// we may encounter either a closing
		// <a> tag or another opening <a> tag
	}

	return result, nil
}

// removeByte takes a slice of bytes and removes all occurences of a given byte from it.
func removeByte(buffer []byte, value byte) []byte {
	res := make([]byte, 0, len(buffer))
	// Add all previous values to the new slice but the exception
	for _, v := range buffer {
		if v != value {
			res = append(res, v)
		}
	}

	return res
}
