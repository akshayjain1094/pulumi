// Copyright 2016 Marapongo, Inc. All rights reserved.

package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAsName(t *testing.T) {
	var goodNames = []string{
		"simple",                              // all alpha.
		"SiMplE",                              // mixed-case alpha.
		"simple0",                             // alphanumeric.
		"SiMpLe0",                             // mixed-case alphanumeric.
		"_",                                   // permit underscore.
		"s1MPl3_",                             // mixed-case alphanumeric/underscore.
		"_s1MPl3",                             // ditto.
		"namespace/complex",                   // multi-part name.
		"_naMeSpace0/coMpl3x32",               // multi-part, alphanumeric, etc. name.
		"n_ameSpace3/moRenam3sp4ce/_Complex5", // even more complex parts.
	}
	for _, nm := range goodNames {
		assert.Equal(t, true, IsName(nm), "IsName expected to be true: %v", nm)
		assert.Equal(t, nm, string(AsName(nm)), "AsName expected to echo back: %v", nm)
	}

	var badNames = []string{
		"0_s1MPl3",                         // cannot start with a number.
		"namespace/0complex",               // ditto.
		"namespace/morenamespace/0complex", // ditto.
		"s!mple",                          // bad characters.
		"namesp@ce/complex",               // ditto.
		"namespace/morenamespace/compl#x", // ditto.
	}
	for _, nm := range badNames {
		assert.Equal(t, false, IsName(nm), "IsName expected to be false: %v", nm)
	}
}

func TestNameSimple(t *testing.T) {
	assert.Equal(t, "simple", string(AsName("simple").Simple()))
	assert.Equal(t, "complex", string(AsName("namespace/complex").Simple()))
	assert.Equal(t, "complex", string(AsName("ns1/ns2/ns3/ns4/complex").Simple()))
	assert.Equal(t, "c0Mpl3x_", string(AsName("_/_/_/_/a0/c0Mpl3x_").Simple()))
}

func TestNameNamespace(t *testing.T) {
	assert.Equal(t, "", string(AsName("simple").Namespace()))
	assert.Equal(t, "namespace", string(AsName("namespace/complex").Namespace()))
	assert.Equal(t, "ns1/ns2/ns3/ns4", string(AsName("ns1/ns2/ns3/ns4/complex").Namespace()))
	assert.Equal(t, "_/_/_/_/a0", string(AsName("_/_/_/_/a0/c0Mpl3x_").Namespace()))
}

func TestRefParse(t *testing.T) {
	{
		p := Ref("simple").parse()
		assert.Equal(t, "", p.Proto)
		assert.Equal(t, DefaultRefBase, p.Base)
		assert.Equal(t, "simple", string(p.Name))
	}
	{
		p := Ref("namespace/complex").parse()
		assert.Equal(t, "", p.Proto)
		assert.Equal(t, DefaultRefBase, p.Base)
		assert.Equal(t, "namespace/complex", string(p.Name))
	}
	{
		p := Ref("ns1/ns2/ns3/ns4/complex").parse()
		assert.Equal(t, "", p.Proto)
		assert.Equal(t, DefaultRefBase, p.Base)
		assert.Equal(t, "ns1/ns2/ns3/ns4/complex", string(p.Name))
	}
	{
		p := Ref("_/_/_/_/a0/c0Mpl3x_").parse()
		assert.Equal(t, "", p.Proto)
		assert.Equal(t, DefaultRefBase, p.Base)
		assert.Equal(t, "_/_/_/_/a0/c0Mpl3x_", string(p.Name))
	}
	{
		p := Ref("github.com/ns1/ns2/ns3/ns4/complex").parse()
		assert.Equal(t, "", p.Proto)
		assert.Equal(t, "github.com/", p.Base)
		assert.Equal(t, "ns1/ns2/ns3/ns4/complex", string(p.Name))
	}
	{
		p := Ref("git://github.com/ns1/ns2/ns3/ns4/complex").parse()
		assert.Equal(t, "git://", p.Proto)
		assert.Equal(t, "github.com/", p.Base)
		assert.Equal(t, "ns1/ns2/ns3/ns4/complex", string(p.Name))
	}
	{
		r := Ref("git://github.com/ns1/ns2/ns3/ns4/complex")
		p := r.parse()
		assert.Equal(t, p.Proto, r.Proto())
		assert.Equal(t, p.Base, r.Base())
		assert.Equal(t, p.Name, r.Name())
	}
	{
		r := Ref("git://github.com/_/_/_/_/a0/c0Mpl3x_")
		p := r.parse()
		assert.Equal(t, p.Proto, r.Proto())
		assert.Equal(t, p.Base, r.Base())
		assert.Equal(t, p.Name, r.Name())
	}
}