import CustomAxios from '../../../util/axiosUtil'
import config from '../../config'
class TurretApi {
    prefix: string;

    constructor() {
        this.prefix = '';
    }

    async command(data) {
        const response = await CustomAxios.post(`${config.backendHost}${this.prefix}/RotaryCommand`, data)
        return response.data;
    }

    async commandAcknowledge(moduleType) {
        const response = await CustomAxios.get(`${config.backendHost}${this.prefix}/RotaryCommandAcknowledge?moduleType=${moduleType}`,)
        return response.data;
    }



    // async getStatus() {
    //     const response = await CustomAxios.get(`${config.backendHost}/${this.prefix}/status`)
    //     return response.data;
    // }

    // async setLaser(checked) {
    //     const response = await CustomAxios.post(`${config.backendHost}/${this.prefix}/setLaser`, {
    //         checked: checked
    //     })
    //     return response.data;
    // }

    // async start() {
    //     const response = await CustomAxios.post(`${config.backendHost}/${this.prefix}/start`)
    //     return response.data;
    // }

    // async setAngle(values) {
    //     const response = await CustomAxios.post(`${config.backendHost}/${this.prefix}/setAngle`, values)
    //     return response.data;
    // }

    // async setPartAngle(values) {
    //     const response = await CustomAxios.post(`${config.backendHost}/${this.prefix}/setPartAngle`, values)
    //     return response.data;
    // }


    // async setSpeed(values) {
    //     const response = await CustomAxios.post(`${config.backendHost}/${this.prefix}/setSpeed`, values)
    //     return response.data;
    // }

    // async stop() {
    //     const response = await CustomAxios.get(`${config.backendHost}/${this.prefix}/stop`)
    //     return response.data;
    // }
}
const turretApi = new TurretApi()
export default turretApi