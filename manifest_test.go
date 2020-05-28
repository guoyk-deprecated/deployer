package main

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

const (
	testSpec = `
cpu: 50:250
build:
  - value 1
  - value 2
package:
  - value 3
  - value 4

test:
  mem: 64:256
  build:
    - t value 1
    - t value 2
uat:
  package:
    - s value 3
    - s value 4
`

	testSpec2 = `
default:
  cpu: 50:250
  build:
    - value {{__PROFILE__}}
    - value 2
  package:
    - value 3
    - value 4
    - value {{__key1__uppercase__}}

test:
  mem: 64:256
  vars:
    key1: val1
  build:
    - t value {{__profile__}}
    - t value 2
uat:
  package:
    - s value 3
    - s value 4
`
)

func TestManifest_Profile(t *testing.T) {
	s := Manifest{}
	err := yaml.Unmarshal([]byte(testSpec), &s)
	require.NoError(t, err)
	require.Equal(t, []string{"s value 3", "s value 4"}, s.Profile("uat").Package)
	require.Equal(t, []string{"value 1", "value 2"}, s.Profile("uat").Build)
	require.Equal(t, []string{"value 3", "value 4"}, s.Profile("test").Package)
	require.Equal(t, []string{"t value 1", "t value 2"}, s.Profile("test").Build)
	require.Equal(t, []byte("#!/bin/bash\nset -eux\nt value 1\nt value 2\n"), s.Profile("test").GenerateBuild())
	require.Equal(t, []byte("s value 3\ns value 4\n"), s.Profile("uat").GenerateDockerfile())
	require.Equal(t, "50m", s.Profile("test").RequestsCPU())
	require.Equal(t, "250m", s.Profile("test").LimitsCPU())
	require.Equal(t, "64Mi", s.Profile("test").RequestsMEM())
	require.Equal(t, "256Mi", s.Profile("test").LimitsMEM())
	require.Equal(t, "", s.Profile("uat").LimitsMEM())

	s = Manifest{}
	err = yaml.Unmarshal([]byte(testSpec2), &s)
	require.NoError(t, err)
	require.Equal(t, []string{"s value 3", "s value 4"}, s.Profile("uat").Package)
	require.Equal(t, []string{"value UAT", "value 2"}, s.Profile("uat").Build)
	require.Equal(t, []string{"value 3", "value 4", "value VAL1"}, s.Profile("test").Package)
	require.Equal(t, []string{"t value test", "t value 2"}, s.Profile("test").Build)
	require.Equal(t, []byte("#!/bin/bash\nset -eux\nt value test\nt value 2\n"), s.Profile("test").GenerateBuild())
	require.Equal(t, []byte("s value 3\ns value 4\n"), s.Profile("uat").GenerateDockerfile())
	require.Equal(t, "50m", s.Profile("test").RequestsCPU())
	require.Equal(t, "250m", s.Profile("test").LimitsCPU())
	require.Equal(t, "64Mi", s.Profile("test").RequestsMEM())
	require.Equal(t, "256Mi", s.Profile("test").LimitsMEM())
	require.Equal(t, "", s.Profile("uat").LimitsMEM())
}
