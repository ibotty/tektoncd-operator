//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package remote

import (
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pkg/errors"
	"github.com/sigstore/cosign/internal/oci"
)

// SignedImageIndex provides access to a remote index reference, and its signatures.
func SignedImageIndex(ref name.Reference, options ...Option) (oci.SignedImageIndex, error) {
	o, err := makeOptions(ref.Context(), options...)
	if err != nil {
		return nil, err
	}
	ri, err := remoteIndex(ref, o.ROpt...)
	if err != nil {
		return nil, err
	}
	return &index{
		v1Index: ri,
		ref:     ref,
		opt:     o,
	}, nil
}

// We alias ImageIndex so that we can inline it without the type
// name colliding with the name of a method it had to implement.
type v1Index v1.ImageIndex

type index struct {
	v1Index
	ref name.Reference
	opt *options
}

var _ oci.SignedImageIndex = (*index)(nil)

// Signatures implements oic.SignedImageIndex
func (i *index) Signatures() (oci.Signatures, error) {
	return signatures(i, i.opt)
}

// Attestations implements oic.SignedImageIndex
func (i *index) Attestations() (oci.Attestations, error) {
	// TODO(mattmoor): allow accessing attestations.
	return nil, errors.New("NYI")
}

// SignedImage implements oic.SignedImageIndex
func (i *index) SignedImage(h v1.Hash) (oci.SignedImage, error) {
	img, err := i.Image(h)
	if err != nil {
		return nil, err
	}
	return &image{
		Image: img,
		opt:   i.opt,
	}, nil
}

// SignedImageIndex implements oic.SignedImageIndex
func (i *index) SignedImageIndex(h v1.Hash) (oci.SignedImageIndex, error) {
	ii, err := i.ImageIndex(h)
	if err != nil {
		return nil, err
	}
	return &index{
		v1Index: ii,
		opt:     i.opt,
	}, nil
}
