import Vue from "vue";

Vue.prototype.$deepCopyObject = (entrada) => {
    let isArray = false;
    if (Array.isArray(entrada)) {
        isArray = true;
    }
    if (typeof entrada != "object" && !isArray) {
        console.error("not an object!");
        return;
    }
    let salida = Object.assign({}, entrada);
    for (var elemento in entrada) {
        if (Object.prototype.hasOwnProperty.call(entrada, elemento)) {
            var attr = entrada[elemento];
            if (typeof attr === "object") {
                //console.log("have an object", attr, elemento);
                salida[elemento] = Vue.prototype.$deepCopyObject(attr);
            } else if (Array.isArray(attr)) {
                //console.log("have an array", attr, elemento);
                let arreglo = Object.assign([], attr);
                for (let i = 0; i < arreglo.length; i++) {
                    if (typeof arreglo[i] == "object") {
                        arreglo[i] = Vue.prototype.$deepCopyObject(arreglo[i]);
                    }
                }
                salida[elemento] = arreglo;
            }
        }
    }
    if (isArray) {
        salida = Object.values(salida);
    }
    return salida;
};