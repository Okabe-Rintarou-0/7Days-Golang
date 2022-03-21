import React from "react";
import {Card, Col, Row} from 'antd'
import EChartsReact from "echarts-for-react";
import {get} from "../utils/ajax";

interface MonitorProps {
    addr: string
    namespace: string
}

interface Info {
    maxVolume: number
    used: number
    percent: number
}

export default class Monitor extends React.Component<MonitorProps, any> {
    state = {
        addr: "unknown",
        namespace: "unknown",
        sampleInterval: 500,
        maxVolume: 0,
        percent: Array<number>(),
        rtt: Array<number>(),
    };

    constructor(props: MonitorProps) {
        super(props);
        this.setState({
            addr: props.addr,
            namespace: props.namespace,
        });
    }

    componentDidMount(): void {
        this.poll();
    }

    reset = () => {
        this.setState({
            percent: [],
            rtt: [],
        });
    };

    getUsageOption = () => {
        return {
            title: {
                text: 'Cash usage',
                subtext: `Group: ${this.props.namespace}`
            },
            toolbox: {
                right: 100,
                show: true,
                feature: {
                    // magicType: {type: ['line',]},
                    // restore: {},
                    saveAsImage: {}
                }
            },
            legend: {
                data: ['percent']
            },
            xAxis: {
                type: 'category',
                boundaryGap: false,
                data: []
            },
            yAxis: {
                type: 'value',
                max: 100,
                axisLabel: {
                    formatter: '{value}%'
                }
            },
            series: [
                {
                    name: 'percent',
                    type: 'line',
                    data: this.state.percent,
                    areaStyle: {
                        normal: {
                            color: '#275F82' //改变区域颜色
                        }
                    },
                    itemStyle: {
                        normal: {
                            color: '#275F82', //改变折线点的颜色
                            lineStyle: {
                                color: '#253A5D' //改变折线颜色
                            }
                        }
                    },
                },
            ]
        };
    };

    getRTTOption = () => {
        return {
            title: {
                text: 'RTT',
                subtext: `Group: ${this.props.namespace}`
            },
            toolbox: {
                right: 100,
                show: true,
                feature: {
                    // magicType: {type: ['line',]},
                    // restore: {},
                    saveAsImage: {}
                }
            },
            legend: {
                data: ['rtt']
            },
            xAxis: {
                type: 'category',
                boundaryGap: false,
                data: []
            },
            yAxis: {
                type: 'value',
                axisLabel: {
                    formatter: '{value}ms'
                }
            },
            series: [
                {
                    name: 'rtt',
                    type: 'line',
                    data: this.state.rtt,
                    areaStyle: {
                        normal: {
                            color: '#87b541' //改变区域颜色
                        }
                    },
                    itemStyle: {
                        normal: {
                            color: '#118218', //改变折线点的颜色
                            lineStyle: {
                                color: '#1f5d1a' //改变折线颜色
                            }
                        }
                    },
                }
            ]
        };
    };

    defaultPath = '/__cash__/';

    pushBack(array: Array<number>, data: number, maxLength: number): Array<number> {
        if (array.length === maxLength) {
            array = array.slice(1)
        }
        return [...array, data]
    }

    getInfo = () => {
        if (this.props.namespace.length == 0) return;
        let url: string = `http://${this.props.addr}${this.defaultPath}${this.props.namespace}/info`;
        let start = Date.now();
        get(url, (info: Info) => {
            this.setState({
                percent: this.pushBack(this.state.percent, Number(info.percent), 50),
                rtt: this.pushBack(this.state.rtt, Date.now() - start, 50)
            })
        }, (err: any) => {
        })
    };

    poll = () => setInterval(() => {
        this.getInfo();
    }, this.state.sampleInterval);

    charts = () => (
        <div>
            <EChartsReact option={this.getUsageOption()}
                          style={{height: '500px'}}/>
            <EChartsReact option={this.getRTTOption()}
                          style={{height: '500px'}}/>
        </div>
    );

    average = (arr: Array<number>): number => {
        if (arr.length == 0) {
            return NaN
        } else {
            return this.formatFloat(arr.reduce(((previousValue: number, currentValue: number): number => previousValue + currentValue)) / arr.length)
        }
    };

    formatFloat = (flt: number) => {
        return Math.floor(flt * 100 + 0.5) / 100
    };

    render()
        :
        React.ReactElement<any, string | React.JSXElementConstructor<any>> | string | number | {} | Iterable<React.ReactNode> | React.ReactPortal | boolean | null | undefined {
        return <Row>
            <Col span={12}> {this.charts()}</Col>
            <Col span={12}>
                <div className="site-card-border-less-wrapper">
                    <Card title="信息一览" bordered={false} headStyle={{fontSize: "2em"}}
                          bodyStyle={{fontSize: "1.5em"}}
                          style={{textAlign: "left"}}>
                        <p>地址: {this.props.addr}</p>
                        <p>命名空间: {this.props.namespace}</p>
                        <p>平均空间占用率: {this.average(this.state.percent)} %</p>
                        <p>平均响应时间: {this.average(this.state.rtt)} ms</p>
                    </Card>
                </div>
            </Col>
        </Row>
    }
}