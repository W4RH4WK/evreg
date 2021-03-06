<!DOCTYPE html>
<html>

<head>
    <title>Event Registry - Register {{.Name}}</title>
    <link rel="icon" type="image/png" href="../../../assets/icon.png" />
    <link rel="stylesheet" href="../../../assets/bootstrap.min.css" />
</head>

<body>
    <div class="container">
        <h1>Register for Slot {{.Name}}</h1>

        <form method="POST">
            <div class="form-group">
                <label for="name">Student ID Number</label>
                <input type="text" class="form-control" name="name" pattern="[0-9]+" />
            </div>
            <button type="submit" class="btn btn-primary">Submit</button>
        </form>

        <hr />
        <div class="text-center">
            <img src="../../../assets/logo.png" title="Logo" alt="Logo" />
            <br />
            <a href="../../../assets/license.txt">License</a>
        </div>
    </div>
</body>

</html>
