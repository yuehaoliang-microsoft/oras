/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package manifest

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras/cmd/oras/internal/file"
	"oras.land/oras/cmd/oras/internal/option"
)

type pushOptions struct {
	option.Common
	option.Descriptor
	option.Pretty
	option.Remote

	targetRef string
	fileRef   string
	mediaType string
}

func pushCmd() *cobra.Command {
	var opts pushOptions
	cmd := &cobra.Command{
		Use:   "push name[:tag|@digest] file",
		Short: "[Preview] Push a manifest to remote registry",
		Long: `[Preview] Push a manifest to remote registry

** This command is in preview and under development. **

Example - Push a manifest to repository 'locahost:5000/hello' and tag with 'latest':
  oras manifest push localhost:5000/hello:latest manifest.json

Example - Push a manifest to repository 'locahost:5000/hello' and output the prettified descriptor:
oras manifest push --descriptor --pretty localhost:5000/hello manifest.json

Example - Push an ORAS artifact manifest to repository 'locahost:5000/hello' and tag with 'latest':
  oras manifest push --media-type application/vnd.cncf.oras.artifact.manifest.v1+json localhost:5000/hello:latest oras_manifest.json

Example - Push a manifest without TLS:
  oras manifest push --insecure localhost:5000/hello:latest manifest.json
`,
		Args: cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.ReadPassword()
		},
		RunE: func(_ *cobra.Command, args []string) error {
			opts.targetRef = args[0]
			opts.fileRef = args[1]
			return pushManifest(opts)
		},
	}

	option.ApplyFlags(&opts, cmd.Flags())
	cmd.Flags().StringVarP(&opts.mediaType, "media-type", "", "", "media type of manifest")
	return cmd
}

func pushManifest(opts pushOptions) error {
	ctx, _ := opts.SetLoggerLevel()
	repo, err := opts.NewRepository(opts.targetRef, opts.Common)
	if err != nil {
		return err
	}

	var mediaType string
	if opts.mediaType != "" {
		mediaType = opts.mediaType
	} else {
		mediaType, err = file.ParseMediaType(opts.fileRef)
		if err != nil {
			return err
		}
	}

	// prepare manifest content
	desc, rc, err := file.PrepareContent(opts.fileRef, mediaType, "", 0)
	if err != nil {
		return err
	}
	defer rc.Close()

	var ref string
	if tag := repo.Reference.Reference; tag != "" {
		ref = tag
	} else {
		ref = desc.Digest.String()
	}

	got, err := repo.Resolve(ctx, ref)
	if errors.Is(err, errdef.ErrNotFound) ||
		(err == nil && got.Digest != desc.Digest) {
		err = repo.PushReference(ctx, desc, rc, ref)
	}
	if err != nil {
		return err
	}

	// outputs manifest's descriptor
	if opts.OutputDescriptor {
		desc.MediaType = opts.mediaType
		bytes, err := opts.Marshal(desc)
		if err != nil {
			return err
		}
		return opts.Output(os.Stdout, bytes)
	}

	fmt.Println("Pushed", opts.targetRef)
	fmt.Println("Digest:", desc.Digest)

	return nil
}
