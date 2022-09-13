import CustomAxios from '../../../util/axiosUtil'
import config from '../../config'
class ConfigApi {
    prefix: string;

    constructor() {
        this.prefix = '';
        // this.prefix = '/radar/fc2/config';
    }

    async set(data: any) {
        return new Promise<any>((resolve, reject) => {
            CustomAxios.post(`${config.backendHost}${this.prefix}/ConfigFc2Param`, data).then((response) => {
                resolve(response)
            }).catch((err) => {
                reject(err)
            });

        })

    }

    // async ack() {
    //     const response = await CustomAxios.get(`${config.backendHost}${this.prefix}/Acknowledged`)
    //     return response.data;
    // }
}
const turretApi = new ConfigApi()
export default turretApi