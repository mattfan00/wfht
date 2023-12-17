document.addEventListener("alpine:init", () => {
    Alpine.data("calendar", () => ({
        selected: [],
        type: 0,
        isSelected(d) {
            return this.selected.includes(d)
        },
        addSelected(date) {
            if (!this.isSelected(date)) {
                this.selected.push(date)
            } else {
                this.selected.splice(this.selected.indexOf(date), 1)
            }
        },
        eventStyle(isCheckIn, date) {
            style = ""
            if (this.isSelected(date)) {
                style = `${style} bg-black text-white`
            } else if (isCheckIn) {
                style = `${style} bg-green-300`
            }

            return style
        },
        submitSelected(type) {
            htmx.ajax('POST', '/events', {
                values: {
                    dates: this.selected, 
                    type: type
                }
            })
        }
    }))
})

function formatDate(d) {
    function pad(v) {
        return ("0" + v).slice(-2)
    }

    const day = pad(d.getDate())
    const month = pad(d.getMonth() + 1)
    const year = d.getFullYear()

    return `${year}-${month}-${day}`
}


