/*
MIT License

# Copyright (c) 2019 Andrzej Rehmann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package xml

import (
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

type Credential struct {
	Tags map[string]string
}

/*
Converts credentials.xml into a slice of structs with all fields reduced.
XML version is ignored as I could no find a parser which could handle xml 1.0+
Jenkins credentials.xml is using xml 1.1 but it does not seem to be using any of the new features.
With xml 1.0+ this can eventually blow up.
*/
func ParseCredentialsXml(credentialsXml []byte) (*[]Credential, error) {
	credentialsXpaths := []string{"//java.util.concurrent.CopyOnWriteArrayList/*", "//list/*"}
	credentials := make([]Credential, 0)
	credentialsDocument, err := parseXml(credentialsXml)
	if err != nil {
		return &credentials, err
	}
	for _, credentialsXpath := range credentialsXpaths {
		for _, credentialNode := range credentialsDocument.FindElements(credentialsXpath) {
			credential := &Credential{
				Tags: map[string]string{},
			}
			for _, child := range credentialNode.ChildElements() {
				reduceFields(child, credential)
			}
			credentials = append(credentials, *credential)
		}
	}
	return &credentials, nil
}

/*
There is a possibility that a field could get overridden but I haven't seen an example of that yet.
*/
func reduceFields(node *etree.Element, credential *Credential) {
	credential.Tags[node.Tag] = strings.TrimSpace(node.Text())
	for _, child := range node.ChildElements() {
		credential.Tags[child.Tag] = strings.TrimSpace(child.Text())
		reduceFields(child, credential)
	}
}

func parseXml(credentialsXml []byte) (*etree.Document, error) {
	document := etree.NewDocument()
	err := document.ReadFromString(stripXmlVersion(credentialsXml))
	if err != nil {
		return nil, err
	}
	return document, nil
}

/*
HACK ALERT:
Stripping xml version because I could not find any decoder which would parse xml version 1.0+
Jenkins uses xml version 1.1+ so this may blow up.
*/
func stripXmlVersion(credentials []byte) string {
	return regexp.
		MustCompile("(?m)^.*<?xml.*$").
		ReplaceAllString(string(credentials), "")
}
