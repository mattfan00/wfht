function formatDate(d) {
    function pad(v) {
        return ("0" + v).slice(-2)
    }

    const day = pad(d.getDate())
    const month = pad(d.getMonth() + 1)
    const year = d.getFullYear()

    return `${year}-${month}-${day}`
}
