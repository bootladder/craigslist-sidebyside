let Appbar = mui.react.Appbar;
let Button = mui.react.Button;
let Container = mui.react.Container;
let Row = mui.react.Row;
let Col = mui.react.Col;
let Input = mui.react.Input;

class CraigslistQueryColumnUserInput extends React.Component {
    constructor(props) {
        super(props);
        this.loadButtonClicked = this.loadButtonClicked.bind(this);
        this.deleteButtonClicked = this.deleteButtonClicked.bind(this);

        this.categorySelectorSelected = this.categorySelectorSelected.bind(this);
        this.citySelectorSelected = this.citySelectorSelected.bind(this);
        this.cityInputEvent = this.cityInputEvent.bind(this);
        this.queryInputEvent = this.queryInputEvent.bind(this);

        this.state = {
            category: "category1",
            city: "city1",
            query: "query1"
        }
    }
    loadButtonClicked(){
        console.log("loadButtonClicked: doRequest: columnIndex: "+this.props.columnIndex + "URL: "+this.URLInput.controlEl.value)
        this.props.doRequest(this.props.columnIndex,this.URLInput.controlEl.value)
    }
    deleteButtonClicked(){
        console.log("deletebuttonclicked: columnIndex:  " + this.props.columnIndex)
        this.props.doDeleteRequest(this.props.columnIndex)
    }
///////////////////////////////////
    categorySelectorSelected(e){
        console.log("categorySelector Selected" + JSON.stringify(e))
        this.setState({ category: e })
        console.log("state is:  " + JSON.stringify(this.state))
        this.URLInput.controlEl.value = this.createCraigslistURL(this.state.city,e,this.state.query)
    }
    citySelectorSelected(e){
        console.log("citySelector Selected" + JSON.stringify(e))
        this.setState({ city: e })
        this.URLInput.controlEl.value = this.createCraigslistURL(e,this.state.category,this.state.query)
    }
    cityInputEvent(value){
        console.log("cityInputEvent " + value)
        this.setState({ city: value  })
        this.URLInput.controlEl.value = this.createCraigslistURL(value,this.state.category,this.queryInput.controlEl.value)
    }
    queryInputEvent(e){
        console.log("queryInputEvent " + e.target.value)
        this.setState({ query: this.queryInput.controlEl.value })
        this.URLInput.controlEl.value = this.createCraigslistURL(this.state.city,this.state.category,this.queryInput.controlEl.value)
    }
    createCraigslistURL(city,category,query){
        return "http://"+city+".craigslist.org/search/"+category+"?sort=date&query="+query
    }
/////////////////////////////////

    componentDidMount(){
        console.log("UserInput componentDidMount() : column doRequest with: "+this.props.url)
        this.props.doRequest(this.props.columnIndex,this.props.url)
        this.URLInput.controlEl.value = this.props.url
    }
/////////////////////////////////

    render() {
        console.log("     render CraigslistQueryUserInput: " + JSON.stringify(this.props))
        return (
            <div>
                <Container>
                <Row>
                    <Input className="mui--text-caption input-100percent" 
                            ref={el => { this.URLInput = el; }} 
                            placeholder="Craigslist Search URL" />
                </Row>
                <Row>
                    <Input className="input-100percent" placeholder="Search Query"
                            ref={el => { this.queryInput = el; }} 
                            onChange={this.queryInputEvent} />
                </Row>
                <Row>
                    <CategorySelector categorySelectorSelected={this.categorySelectorSelected} />
                    <CitySelector citySelectorSelected={this.citySelectorSelected} cityInputEvent={this.cityInputEvent}/>
                </Row>
                <Row>
                    <Button onClick={this.loadButtonClicked} size="small" color="primary">Load Results and Save URL</Button>
                    <Button onClick={this.deleteButtonClicked} size="small" color="primary">Delete this Column</Button>
                </Row>
                </Container>
            </div>
        );
    }
}

class CraigslistQueryColumnResults extends React.Component {

    render() {
        console.log("     render CraigslistQueryColumnResults is: " + this.props.results)
        const html = $.parseHTML( this.props.results )
        console.log("     HTML: "  + JSON.stringify(html))
        
        var resultRows = []
        var i=0
        $(html).find(".result-row").each( function(q) {
            i++
            resultRows.push(React.createElement(
                'div', {key: i}, <div dangerouslySetInnerHTML= {{__html: this.innerHTML}} />
            ))
        })

        return (
            <div>
                <div className="hello x-scroll y-scroll">{resultRows}</div>
            </div>
        );
    }
}

class CraigslistQueryColumn extends React.Component {

    componentDidUpdate(){
        console.log("Top Column componentWilLReceiveProps()")
    }

    render() {
        return (
            <div>
                <div className="mui--text-left">
                    <CraigslistQueryColumnUserInput 
                        doRequest={this.props.doRequest}
                        doDeleteRequest={this.props.doDeleteRequest}
                        url={this.props.url}
                        columnIndex={this.props.columnIndex}
                        />
                </div>
                <CraigslistQueryColumnResults results={this.props.response} />
            </div>
        );
    }
}