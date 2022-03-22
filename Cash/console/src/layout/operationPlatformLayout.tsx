import React, {createRef} from "react";
import {Button, Col, Input, InputRef, Layout, message, Row, Select, Table, Tag} from "antd";
import {Option} from "antd/es/mentions";
import {SendOutlined, SyncOutlined} from '@ant-design/icons'
import {del, getForText, put} from "../utils/ajax";
import format from 'date-fns/format';

require("../css/console.css");

interface RequestEntry {
    key: number
    method: string
    address: string
    group: string
    requestKey: string
    value?: string
    time: string
    status: string
    response: string
    responseTime: number
}

interface InputParams {
    method: string
    address: string
    group: string
    key: string
    value: string | undefined
}

const columns = [
    {title: 'Time', dataIndex: 'time', key: 'time'},
    {
        title: 'Method', dataIndex: 'method', key: 'method', render: (method: string) => {
            let color = 'white';
            switch (method) {
                case "GET":
                    color = "green";
                    break;
                case "DELETE":
                    color = "red";
                    break;
                case "PUT":
                    color = "blue";
                    break;
            }
            return <span><Tag color={color} key={method}>{method}</Tag></span>
        }
    },
    {title: 'Address', dataIndex: 'address', key: 'address'},
    {title: 'Group', dataIndex: 'group', key: 'group'},
    {title: 'Key', dataIndex: 'requestKey', key: 'requestKey'},
    {title: 'Value', dataIndex: 'value', key: 'value'},
    {title: 'Action', key: 'action', render: () => <Button><SyncOutlined/></Button>},
];

const {Header, Content, Footer} = Layout;
const httpMethods = ["GET", "PUT", "DELETE"];
export default class OperationPlatformLayout extends React.Component<any, any> {
    state = {
        selectedMethod: "GET",
        requestEntries: Array<RequestEntry>(),
    };

    private readonly addressInputRef: React.RefObject<InputRef>;
    private readonly keyInputRef: React.RefObject<InputRef>;
    private readonly valueInputRef: React.RefObject<InputRef>;
    private readonly groupInputRef: React.RefObject<InputRef>;

    constructor(props: any) {
        super(props);
        this.addressInputRef = createRef();
        this.keyInputRef = createRef();
        this.valueInputRef = createRef();
        this.groupInputRef = createRef();
    }

    selectMethod = (selectedMethod: string) => {
        this.setState({
            selectedMethod: selectedMethod
        });
    };

    valueInput = () => {
        return this.state.selectedMethod == "PUT" ?
            (<Col>
                <Input ref={this.valueInputRef} addonBefore="value:" placeholder="value"/>
            </Col>) : null
    };

    getInputRefValue = (ref: InputRef | null): string => {
        if (ref !== null && ref.input !== null)
            return ref.input.value;
        return ""
    };

    getInputParams = (): InputParams => {
        return {
            method: this.state.selectedMethod,
            address: this.getInputRefValue(this.addressInputRef.current),
            group: this.getInputRefValue(this.groupInputRef.current),
            key: this.getInputRefValue(this.keyInputRef.current),
            value: this.getInputRefValue(this.valueInputRef.current),
        };
    };

    httpRequest = (method: string, url: string, callback: Function, onError?: Function) => {
        switch (method) {
            case "GET":
                getForText(url, callback, onError);
                break;
            case "PUT":
                put(url, callback, onError);
                break;
            case "DELETE":
                del(url, callback, onError);
        }
    };

    sendRequest = () => {
        let inputParams = this.getInputParams();
        let url = this.formUrl(inputParams.method, inputParams.address, inputParams.group, inputParams.key, inputParams.value);
        let start = Date.now();
        let requestEntry: RequestEntry;

        message.success('Send request successfully.').then(() => {
        });

        this.httpRequest(inputParams.method, url, (response: string) => {
            requestEntry = {
                response: response,
                responseTime: Date.now() - start,
                status: "Ok",
                time: format(start, 'yyyy-mm-dd HH:mm:ss'),
                address: inputParams.address,
                group: inputParams.group,
                requestKey: inputParams.key,
                value: inputParams.value,
                key: this.state.requestEntries.length,
                method: inputParams.method
            };
            this.setState({
                requestEntries: [...this.state.requestEntries, requestEntry]
            });
        }, (err: string) => {
            requestEntry = {
                response: err,
                responseTime: Date.now() - start,
                status: "Error",
                time: format(start, 'yyyy-mm-dd HH:mm:ss'),
                address: inputParams.address,
                group: inputParams.group,
                requestKey: inputParams.key,
                value: inputParams.value,
                key: this.state.requestEntries.length,
                method: inputParams.method
            };
            this.setState({
                requestEntries: [...this.state.requestEntries, requestEntry]
            });
        });
    };

    formUrl = (method: string, addr: string, group: string, key: string, value: string | undefined): string => (
        method !== "PUT" ?
            `http://${addr}/__cash__/${group}?key=${key}` :
            `http://${addr}/__cash__/${group}?key=${key}&value=${value}`
    );

    table = () => {
        return <Table
            columns={columns}
            expandable={{
                expandedRowRender: (entry: RequestEntry) => {
                    return <div>
                        <p>Status: {entry.status}</p>
                        <p>Response Time: {entry.responseTime}</p>
                        <p>Response: {entry.response}</p>
                    </div>
                }
            }}
            pagination={{pageSize: 12}}
            dataSource={this.state.requestEntries}
        />
    };

    selector = () => (
        <Row>
            <Col>
                <Select defaultValue={"GET"}
                        style={{width: 150}} onSelect={this.selectMethod}>
                    {httpMethods.map((method) => <Option key={method}/>)}
                </Select>
            </Col>
            <Col>
                <Input ref={this.addressInputRef} defaultValue="localhost:8000" addonBefore="http://"
                       placeholder="address"/>
            </Col>
            <Col>
                <Input ref={this.groupInputRef} defaultValue="country" addonBefore="group:" placeholder="group"/>
            </Col>
            <Col>
                <Input ref={this.keyInputRef} defaultValue="China" addonBefore="key:" placeholder="key"/>
            </Col>
            {this.valueInput()}
            <Col>
                <Button onClick={this.sendRequest}>Send<SendOutlined/></Button>
            </Col>
        </Row>
    );

    render(): React.ReactElement {
        return <Layout>
            <Header style={{padding: 0}}><p style={{color: "white", fontSize: "2em"}}>Cash</p></Header>
            <Content style={{margin: '24px 16px 0', height: 1125}}>
                {this.table()}
            </Content>
            <Footer style={{textAlign: 'center'}}>
                {this.selector()}
            </Footer>
        </Layout>
    }
}