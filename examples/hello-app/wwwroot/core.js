const go = new Go();

WebAssembly
    .instantiateStreaming(fetch("main.wasm"), go.importObject)
    .then(result => {
        go.run(result.instance);
    })
    .catch((err) => {
        console.error("Error loading WASM:", err);
    });
