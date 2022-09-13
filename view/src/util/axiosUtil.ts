/*
 * @Descripttion: pangoo-dm
 * @version: 1.0
 * @Author: uidq1343
 * @Date: 2021-12-01 14:47:35
 * @LastEditors: uidq1343
 * @LastEditTime: 2022-03-11 11:15:02
 * @content: edit your page content
 */
import axios, { AxiosResponse } from 'axios';
axios.defaults.adapter = require('axios/lib/adapters/http');

import { ICommonResponse } from './responseUtil'

const CustomAxios = axios.create({
    timeout: 3000, // 设置超时时长
})

// export const getAxiosData = (response: AxiosResponse): any => {
//     const data: ICommonResponse = response.data
//     if (data.code !== 0) {
//         throw data.msg;
//     }
//     return data;
// }

function responseOnFulfilled(response: AxiosResponse) {
    const data: ICommonResponse = response.data
    if (data.code !== 0) {
        throw data.msg;
    }
    return data;
}

function responseOnRejected(err) {
    console.log(err)
    if (err.response && err.response.status === 500) {
        return Promise.reject('服务器内部错误。')
    }
    return Promise.reject(err)
}
// 返回后拦截
CustomAxios.interceptors.response.use(responseOnFulfilled, responseOnRejected)

export default CustomAxios