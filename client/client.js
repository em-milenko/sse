
function onClick1() {
    const source = new EventSource("http://localhost:3000/sse?id=1")
    source.addEventListener("update",(event) => {
        console.log("OnMessage Called:")
        console.log(event)
        source.close()
    })
}

function onClick2() {
    const source = new EventSource("http://localhost:3000/sse?id=2")
    source.addEventListener("update",(event) => {
        console.log("OnMessage Called:")
        console.log(event)
        source.close()
    })
}