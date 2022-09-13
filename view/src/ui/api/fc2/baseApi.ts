import CustomAxios from '../../../util/axiosUtil'
import config from '../../config'
class BaseApi {

    async heartBeat() {
        return new Promise<any>((resolve, reject) => {
            CustomAxios.get(`${config.backendHost}/HeartBeatReq`).then((response) => {
                resolve(response)
            }).catch((err) => {
                reject(err)
            });

        })
    }

    async postReq(url,param) {
        return new Promise<any>((resolve, reject) => {
            CustomAxios.post(`${config.backendHost}/${url}`,param).then((response) => {
                resolve(response)
            }).catch((err) => {
                reject(err)
            });

        })
    }
    
    async getReq(url) {
        return new Promise<any>((resolve, reject) => {
            CustomAxios.get(`${config.backendHost}/${url}`).then((response) => {
                resolve(response)
            }).catch((err) => {
                reject(err)
            });

        })
    }


}
const baseApi = new BaseApi()
export default baseApi