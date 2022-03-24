import React, {createRef} from "react";
import {Badge, Button, Col, Input, InputRef, Layout, message, Row, Select, Table, Tag} from "antd";
import {Option} from "antd/es/mentions";
import {PlusOutlined, SendOutlined, SyncOutlined} from '@ant-design/icons'
import {del, getForText, post, put} from "../utils/ajax";
import format from 'date-fns/format';

require("../css/console.css");

interface RequestEntry {
    key: number
    method: string
    address: string
    group: string
    requestKey: string
    url: string
    value?: string
    time: string
    status: "Ok" | "Error"
    response: string
    responseTime: number
}

interface BatchedRequestEntry {
    method: string
    key: string
    value: string | undefined
}

interface BatchedRequest {
    address: string
    group: string
    requests: Array<BatchedRequestEntry>
}

interface BatchedResponse {
    responses: Array<string>
}

interface InputParams {
    method: string
    address: string
    group: string
    key: string
    value: string | undefined
}

interface RequestParamsWrap {
    url: string
    method: string
    address: string
    group: string
    key: string
    value: string | undefined
}

const {Header, Content, Footer} = Layout;
const httpMethods = ["GET", "PUT", "DELETE"];
export default class OperationPlatformLayout extends React.Component<any, any> {
    state = {
        selectedMethod: "GET",
        selectedMode: "Normal Mode",
        requestEntries: Array<RequestEntry>(),
        batchRequestEntries: Array<RequestEntry>(),
    };

    private readonly addressInputRef: React.RefObject<InputRef>;
    private readonly keyInputRef: React.RefObject<InputRef>;
    private readonly valueInputRef: React.RefObject<InputRef>;
    private readonly groupInputRef: React.RefObject<InputRef>;

    mapColor = (method: string): string => {
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
        return color;
    };

    private columns = [
        {title: 'Time', dataIndex: 'time', key: 'time'},
        {
            title: 'Method',
            dataIndex: 'method',
            key: 'method',
            render: (method: string) => <span><Tag color={this.mapColor(method)} key={method}>{method}</Tag></span>
        },
        {title: 'Address', dataIndex: 'address', key: 'address'},
        {title: 'Group', dataIndex: 'group', key: 'group'},
        {title: 'Key', dataIndex: 'requestKey', key: 'requestKey'},
        {title: 'Value', dataIndex: 'value', key: 'value'},
        {
            title: 'Action', key: 'action', render: (entry: RequestEntry) => <Button onClick={() => {
                this.sendRequest({
                    address: entry.address,
                    group: entry.group,
                    key: entry.requestKey,
                    method: entry.method,
                    url: entry.url,
                    value: entry.value
                } as RequestParamsWrap);
            }}><SyncOutlined/></Button>
        },
    ];

    constructor(props: any) {
        super(props);
        this.addressInputRef = createRef();
        this.keyInputRef = createRef();
        this.valueInputRef = createRef();
        this.groupInputRef = createRef();
    }

    onSelectMode = (mode: string) => {
        this.setState({
            selectedMode: mode
        })
    };

    onSelectMethod = (selectedMethod: string) => {
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

    sendRequest = (paramsWrap: RequestParamsWrap) => {
        let start = Date.now();
        let requestEntry: RequestEntry;
        this.httpRequest(paramsWrap.method, paramsWrap.url, (response: string) => {
            let requestEntry = {
                key: this.state.requestEntries.length,
                time: format(start, 'yyyy-mm-dd HH:mm:ss'),
                url: paramsWrap.url,
                method: paramsWrap.method,
                address: paramsWrap.address,
                group: paramsWrap.group,
                requestKey: paramsWrap.key,
                value: paramsWrap.value,
                response: response,
                responseTime: Date.now() - start,
                status: "Ok",
            } as RequestEntry;
            this.setState({
                requestEntries: [...this.state.requestEntries, requestEntry]
            });
        }, (err: string) => {
            requestEntry = {
                key: this.state.requestEntries.length,
                time: format(start, 'yyyy-mm-dd HH:mm:ss'),
                url: paramsWrap.url,
                method: paramsWrap.method,
                address: paramsWrap.address,
                group: paramsWrap.group,
                requestKey: paramsWrap.key,
                value: paramsWrap.value,
                response: err,
                responseTime: Date.now() - start,
                status: "Error",
            } as RequestEntry;
            this.setState({
                requestEntries: [...this.state.requestEntries, requestEntry]
            });
        });
    };

    onSendRequest = () => {
        if (this.state.selectedMode === 'Normal Mode') {
            this.sendNormalRequest();
        } else {
            this.sendBatchRequest();
        }
    };

    sendNormalRequest = () => {
        let inputParams = this.getInputParams();
        let url = this.formUrl(inputParams.method, inputParams.address, inputParams.group, inputParams.key, inputParams.value);

        this.sendRequest({
            address: inputParams.address,
            group: inputParams.group,
            key: inputParams.key,
            method: inputParams.method,
            url: url,
            value: inputParams.value
        } as RequestParamsWrap);

        message.success('Send request successfully.').then(() => {
        });
    };

    formUrl = (method: string, addr: string, group: string, key: string, value: string | undefined): string => (
        method !== "PUT" ?
            `http://${addr}/__cash__/${group}?key=${key}` :
            `http://${addr}/__cash__/${group}?key=${key}&value=${value}`
    );

    status = (status: "Ok" | "Error") => {
        let badge;
        if (status === 'Ok') {
            badge = <Badge status="success"/>;
        } else {
            badge = <Badge status="error"/>
        }
        return <p>Status: {badge}{status}</p>;
    };

    table = () => <Table
        columns={this.columns}
        expandable={{
            expandedRowRender: (entry: RequestEntry) => {
                return <div>
                    {this.status(entry.status)}
                    <p>Response Time: {entry.responseTime}</p>
                    <p>Response: {entry.response}</p>
                </div>
            }
        }}
        pagination={{pageSize: 8}}
        dataSource={this.state.requestEntries}
    />;

    addBatch = () => {
        let inputParams = this.getInputParams();
        let batchRequestEntry: RequestEntry = {
            address: inputParams.address,
            group: inputParams.group,
            key: this.state.batchRequestEntries.length,
            method: inputParams.method,
            requestKey: inputParams.key,
            value: inputParams.value,
            url: this.formUrl(inputParams.method, inputParams.address, inputParams.group, inputParams.key, inputParams.value)
        } as RequestEntry;

        this.setState({
            batchRequestEntries: [...this.state.batchRequestEntries, batchRequestEntry]
        })
    };

    addBatchButton = () => {
        if (this.state.selectedMode === "Normal Mode") return null;
        return <Col>
            <Button onClick={this.addBatch}>Add<PlusOutlined/></Button>
        </Col>
    };

    getBatchedRequest = (): BatchedRequest | null => {
        if (this.state.batchRequestEntries.length === 0) return null;
        let batchedRequests = Array<BatchedRequestEntry>();

        let address = this.state.batchRequestEntries[0].address;
        let group = this.state.batchRequestEntries[0].group;

        this.state.batchRequestEntries.map(entry => {
            batchedRequests.push({
                key: entry.requestKey, method: entry.method, value: entry.value
            } as BatchedRequestEntry);
        });

        return {
            address: address, group: group, requests: batchedRequests
        } as BatchedRequest
    };

    sendBatchRequest = () => {
        let batchRequest = this.getBatchedRequest();
        if (batchRequest == null) return;
        let url = `http://${batchRequest.address}/__cash__/${batchRequest.group}/__batch__`;
        post(url, batchRequest, (response: BatchedResponse) => {
                console.log(response.responses);
            },
            (err: any) => {
                console.log(err);
            });
        this.setState({
            batchRequestEntries: []
        });
    };

    batchTable = () => {
        if (this.state.selectedMode === "Normal Mode") return null;
        const batchColumns = [
            {
                title: 'Method',
                dataIndex: 'method',
                key: 'method',
                render: (method: string) => <span><Tag color={this.mapColor(method)} key={method}>{method}</Tag></span>
            },
            {title: 'Address', dataIndex: 'address', key: 'address'},
            {title: 'Group', dataIndex: 'group', key: 'group'},
            {title: 'Key', dataIndex: 'requestKey', key: 'requestKey'},
            {title: 'Value', dataIndex: 'value', key: 'value'},
        ];
        return <Table
            columns={batchColumns}
            pagination={{pageSize: 5}}
            dataSource={this.state.batchRequestEntries}/>
    };

    selector = () => (
        <Row>
            <Col>
                <Select defaultValue={"Normal Mode"}
                        style={{width: 150}} onSelect={this.onSelectMode}>
                    <Option key={"Normal Mode"}/>
                    <Option key={"Batch Mode"}/>
                </Select>
            </Col>
            <Col>
                <Select defaultValue={"GET"}
                        style={{width: 150}} onSelect={this.onSelectMethod}>
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
            {this.addBatchButton()}
            <Col>
                <Button onClick={this.onSendRequest}>Send<SendOutlined/></Button>
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
                {this.batchTable()}
            </Footer>
        </Layout>
    }
}