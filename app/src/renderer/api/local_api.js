import axios from 'axios'

export const SERVER_ADDR = 'http://127.0.0.1:8042';

export default {
    async api_post(req_type, data) {
        let req_addr = SERVER_ADDR + '/';
        return new Promise(resolve => {
            switch (req_type) {
                case 'game/new':
                    req_addr += `games/new/`
                    break;

                case "game/generate":
                    req_addr += `games/boards/generate/`
                    break;

                case "game/setmode":
                    req_addr += `games/setmode/`
                    break;

                case "game/drawbalot":
                    req_addr += "games/drawbalot/"
                    break;

                case "game/check":
                    req_addr += "games/boards/check/"
                    break;

                default:
                    throw ("invalid req_type " + req_type)
            }

            axios.post(req_addr, data)
                .then(response => {
                    console.log('got resp', response);
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
        console.log('api_get', req_type);
        return new Promise((resolve) => {
            let rObj = {}
            let req_addr = SERVER_ADDR + '/';
            switch (req_type) {
                case 'games':
                    req_addr += `games/`
                    break;
                case 'game':
                    req_addr += `games/?gid=${id}`
                    break;
                default:
                    errmsg = "api_get invalid req_type " + req_type
                    console.error(errmsg);
                    throw (errmsg)
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
            let req_addr = SERVER_ADDR + '/';
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
