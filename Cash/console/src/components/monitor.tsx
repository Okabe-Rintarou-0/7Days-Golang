import React from "react";
import {Card, Col, InputNumber, Row} from 'antd'
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

interface GroupInfo {
    name: string
    value: number
}

export default class Monitor extends React.Component<MonitorProps, any> {
    state = {
        addr: "unknown",
        namespace: "unknown",
        sampleInterval: 2000,
        maxVolume: 0,
        percent: Array<number>(),
        rtt: Array<number>(),
        groupInfos: Array<GroupInfo>()
    };

    private pollTimer: NodeJS.Timer | undefined;

    constructor(props: MonitorProps) {
        super(props);
        this.setState({
            addr: props.addr,
            namespace: props.namespace,
        });
    }

    componentDidMount(): void {
        this.pollTimer = this.poll();
    }

    fetchGroupInfo = (groups: Array<string>) => {
        let groupInfos = Array<GroupInfo>();
        groups.forEach((group: string, _: number) => {
            let url = `http://${this.props.addr}${this.defaultPath}${group}/info`;
            get(url, (info: Info) => {
                groupInfos.push({
                    name: group,
                    value: info.maxVolume,
                });
                this.setState({
                    groupInfos: [this.state.groupInfos, ...groupInfos]
                })
            })
        });
    };

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

    getPieChartOption = () => {
        return {
            tooltip: {
                trigger: 'item'
            },
            legend: {
                top: '5%',
                left: 'center'
            },
            series: [
                {
                    name: 'Max Volume',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    avoidLabelOverlap: false,
                    itemStyle: {
                        borderRadius: 10,
                        borderColor: '#fff',
                        borderWidth: 2
                    },
                    label: {
                        show: false,
                        position: 'center'
                    },
                    emphasis: {
                        label: {
                            show: true,
                            fontSize: '40',
                            fontWeight: 'bold'
                        }
                    },
                    labelLine: {
                        show: false
                    },
                    data: this.state.groupInfos
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
        if (this.props.namespace.length === 0) return;
        let url = `http://${this.props.addr}${this.defaultPath}${this.props.namespace}/info`;
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

    lineCharts = () => (
        <div>
            <EChartsReact option={this.getUsageOption()}
                          style={{height: '500px'}}/>
            <EChartsReact option={this.getRTTOption()}
                          style={{height: '500px'}}/>
        </div>
    );

    adjustSampleInterval = (interval: string) => {
        if (this.pollTimer !== undefined) {
            clearInterval(this.pollTimer);
        }

        this.setState({
            sampleInterval: Number(interval) * 1000
        }, () => {
            this.pollTimer = this.poll();
        });
    };

    infoCard = () => (
        <div className="site-card-border-less-wrapper">
            <Card title="Info Panel" bordered={true}
                  headStyle={{fontSize: "2.2em", fontWeight: "bold"}}
                  bodyStyle={{fontSize: "1.7em"}}
                  style={{textAlign: "left"}}>
                <p>Address: {this.props.addr}</p>
                <p>Namespace: {this.props.namespace}</p>
                <p>Average Usage: {this.average(this.state.percent)} %</p>
                <p>Average RTT: {this.average(this.state.rtt)} ms</p>
                <p>Sampling Interval：
                    <InputNumber<string> style={{width: 150}}
                                         addonAfter="seconds"
                                         defaultValue="2"
                                         min="0.5"
                                         max="100"
                                         step="0.5"
                                         onChange={this.adjustSampleInterval}
                                         stringMode
                    />
                </p>
            </Card>
        </div>
    );

    pieChart = () => (
        <div>
            <EChartsReact option={this.getPieChartOption()}
                          style={{height: '600px'}}/>
        </div>
    );

    average = (arr: Array<number>): number => {
        if (arr.length === 0) {
            return NaN
        } else {
            return this.formatFloat(arr.reduce(((previousValue: number, currentValue: number): number => previousValue + currentValue)) / arr.length)
        }
    };

    formatFloat = (flt: number) => {
        return Math.floor(flt * 100 + 0.5) / 100
    };

    render(): React.ReactElement {
        return <Row>
            <Col span={12}> {this.lineCharts()}</Col>
            <Col span={12}>
                {this.infoCard()}
                {this.pieChart()}
            </Col>
        </Row>
    }
}