import axios from 'axios'

export default {
    static SERVER_ADDR = 'localhost:8042';

    async api_post(req_type, data) {
        let req_addr = this.SERVER_ADDR;
        return new Promise(resolve => {
            switch (req_type) {
                case 'game/new':
                    req_addr += `games/new/`
                    break;

                default:
                    throw ("invalid req_type " + req_type)
            }

            axios.post(req_addr, data)
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.error("api_post: ", error.response.data.error);
                    if (error.response && error.response.data) resolve(error.response.data)
                    else resolve(null);
                })
        });
    },

    async api_get(req_type, id, query) {
        return new Promise((resolve) => {
            if (!id) {
                throw (`api_get: missing id ${req_type}`)
            }
            let rObj = {}
            let req_addr = dirapi
            switch (req_type) {
                case 'game':
                    req_addr += `games/${id}/`
                    break;
                default:
                    throw ("api_get invalid req_type " + req_type)
            }

            if (query && query != '') {
                req_addr += `?${query}`
            }

            axios.get(req_addr)
                .then((response) => {
                    console.log("api_get: got response");
                    rObj = response.data
                    resolve(rObj);
                })
                .catch(e => {
                    console.error("api_get: ", e)
                    resolve(null)
                })
        });
    },

    async api_delete(req_type, id, data = {}) {
        return new Promise(resolve => {
            let req_addr = this.SERVER_ADDR
            switch (tipo) {
                case 'game':
                    req_addr += `/games/${id}/`
                    break;
                default:
                    throw ("api_delete: invalid req_type " + tipo)
            }

            axios.delete(req_addr, {
                data: data
            })
                .then(response => {
                    resolve(response.data);
                })
                .catch(error => {
                    console.error(error)
                    resolve(null);
                })
        });
    },
}
