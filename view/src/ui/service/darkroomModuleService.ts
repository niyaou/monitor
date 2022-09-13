import turretApi from '../api/fc2/turretApi'
import numeral from 'numeral';

const moduleType = 2;
class DarkroomModuleService {

    async getStatus() {
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: 'RED_angle'
        })
        return response;
    }

    async getStatusAck() {
        const response = await turretApi.commandAcknowledge(moduleType)  
        return parseStatus(response.commandPayload);
    }

    async start() {
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: 'SET_start'
        })
        return response;
    }

    async setAngle(values) {
        const angleFomat = '+000.00'
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: `SET_angle_${numeral(values.line1.range).format(angleFomat)}_${numeral(values.line2.range).format(angleFomat)}_${numeral(values.radian1.range).format(angleFomat)}_${numeral(values.radian2.range).format(angleFomat)}_${numeral(values.dividedCircle.range).format(angleFomat)}`,
        })
        return response;
    }

    async setPartAngle(data) {
        const angleFomat = '+000.00'
        let cmd;
        switch (data.name) {
            case 'line1':
                cmd = `SET_11111_${numeral(data.value).format(angleFomat)}`
                break;
            case 'line2':
                cmd = `SET_22222_${numeral(data.value).format(angleFomat)}`
                break;
            case 'radian1':
                cmd = `SET_33333_${numeral(data.value).format(angleFomat)}`
                break;
            case 'radian2':
                cmd = `SET_44444_${numeral(data.value).format(angleFomat)}`
                break;
            case 'dividedCircle':
                cmd = `SET_55555_${numeral(data.value).format(angleFomat)}`
                break;
            default:
                throw (`unkown data.type: ${data.type}`)
        }
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: cmd
        })
        return response;
    }


    async setSpeed(values) {
        const speedFomat = '000.00'
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: `SET_speed_${numeral(values.line1.maxSpeed).format(speedFomat)}_${numeral(values.line2.maxSpeed).format(speedFomat)}_${numeral(values.radian1.maxSpeed).format(speedFomat)}_${numeral(values.radian2.maxSpeed).format(speedFomat)}_${numeral(values.dividedCircle.maxSpeed).format(speedFomat)}`,
        })
        return response;
    }

    async stop() {
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: "SET_cease"
        })
        return response;
    }
}

function parseStatus(message) {
    const data = {
        'workStatus': 'idle',
        'data': {
            "laser": true,
            "line1": {
                "angle": 0,
                "speed": 0
            },
            "line2": {
                "angle": 0,
                "speed": 0
            },
            "radian1": {
                "angle": 0,
                "speed": 0
            },
            "radian2": {
                "angle": 0,
                "speed": 0
            },
            "dividedCircle": {
                "angle": 0,
                "speed": 0
            }
        }
    };

    data['data']['line1']['angle'] = Number(message.substring(10, 17));
    data['data']['line2']['angle'] = Number(message.substring(19, 25));
    data['data']['radian1']['angle'] = Number(message.substring(27, 33));
    data['data']['radian2']['angle'] = Number(message.substring(35, 41));
    data['data']['dividedCircle']['angle'] = Number(message.substring(43, 49));
    if ('RED' === message.substring(0, 3)) {
        data['workStatus'] = 'work';
    } else {
        data['workStatus'] = 'idle';
    }
    return data;
}
const darkroomModuleService = new DarkroomModuleService()
export default darkroomModuleService