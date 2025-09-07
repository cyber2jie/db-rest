function transform(row) {
    for (var i = 0; i < row.columns.length; i++) {
        var column = row.columns[i];
        if (column.value) {
            column.value = column.value.toUpperCase()
        }
    }
}