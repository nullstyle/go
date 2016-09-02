A go library for scientific unit calculations, powering a scientific spreadsheet application.  This library is inspired by ruby-unit and allows for code like:

```go

  // creating values
  x := units.Value{M: "10.001", U: units.Millimeter}
  y := units.Value{M: "10.001", U: units.Millimeter}

  // math
  z, err := x.Add(y)

  // units are composite values, algebraic values
  // acceleration
  Acceleration = units.Div{
    N: units.Meter,
    D: units.Pow{ units.Second, 2 },
  }
   

  // Defining constants
  var (
    EarthGravity = units.MustParseValue("9.807 m/s^2")
  )
```

This spreadsheet application could allow for "input cells" and "output cells" which exposes the sheet as an API to others allowing for collaboration.