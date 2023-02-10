# UE Manifest

an Unreal Engine's manifest parser

## Installation

```bash
go add github.com/gounreal/uemanifest
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/gounreal/uemanifest"
)

func main() {
    var r io.ReadSeeker = ...

	manifest, err := uemanifest.Parse(r)
	handleError(err)

	fmt.Println(manifest.ManifestMeta.AppName)
}
```

## Credit

The initial creation of this project used [JFortniteParse](https://github.com/FabianFG/JFortniteParse/blob/master/src/main/kotlin/me/fungames/jfortniteparse/ue4/manifests/objects) & [egmanifest](https://github.com/er-azh/egmanifest) as a reference, so huge thanks to them.
