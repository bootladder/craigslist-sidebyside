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
        this.queryInputEvent = this.queryInputEvent.bind(this);

        this.state = {
            category: "category1",
            city: "city1",
            query: "query1"
        }
    }
    loadButtonClicked(){
        console.log("the columnIndex is "+this.props.columnIndex)
        this.props.doRequest(this.props.columnIndex,this.URLInput.controlEl.value)
    }
    deleteButtonClicked(){
        console.log("deletebuttonclicked: columnIndex:  " + this.props.columnIndex)
        this.props.doDeleteRequest(this.props.columnIndex)
    }
    categorySelectorSelected(e){
        console.log("categorySelector Selected" + JSON.stringify(e))
        this.setState({
            category: e
        })
        this.URLInput.controlEl.value = this.createCraigslistURL(this.state.city,e,this.state.query)
    }
    citySelectorSelected(e){
        console.log("citySelector Selected" + JSON.stringify(e))
        this.setState({
            city: e
        })
        this.URLInput.controlEl.value = this.createCraigslistURL(e,this.state.category,this.state.query)
    }
    queryInputEvent(e){
        console.log("queryInputEvent " + e.target.value)
        this.setState({
            query: this.queryInput.controlEl.value
        })
        this.URLInput.controlEl.value = this.createCraigslistURL(this.state.city,this.state.category,this.queryInput.controlEl.value)
    }

    createCraigslistURL(city,category,query){
        return "http://"+city+".craigslist.org/search/"+category+"?query="+query
    }

    componentDidMount(){
        this.URLInput.controlEl.value = this.props.url
        this.props.doRequest(this.props.columnIndex,this.props.url)
        console.log("column doRequest with: "+this.props.url)
    }
    render() {
        console.log("render CraigslistQueryUserInput: " + JSON.stringify(this.props))
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
        <CategorySelector categorySelectorSelected={this.categorySelectorSelected} myprop="myprop"/>
        <CitySelector citySelectorSelected={this.citySelectorSelected}/>
    </Row>
    <Row>
        <Col><Input placeholder="City" /></Col>
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
        const html = $.parseHTML( this.props.results )
        
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

    constructor(props) {
        super(props);
        this.doRequest = this.doRequest.bind(this);
    }

    state = {
        columnState: "columnState",
        queryResponseData: "queryResponseData"
    }

    componentDidMount() {
    }

    validateCraigslistURL(url){
        console.log(JSON.stringify(url))
        if(url.length < 5){
            console.log("length too short, defaulting URL length was "+url.length)
            return "https://baltimore.craigslist.org/d/architect-engineer-cad/search/egr";
        }
        else return url;
    }
    doRequest(index, craigslistSearchURL){
        console.log("do request" + craigslistSearchURL + " state is " + JSON.stringify(this.state))
        var validatedURL = this.validateCraigslistURL(craigslistSearchURL)
        console.log("validatedURL: " + this.validatedURL)

        var myJsonRequestObj = {
            searchURL: encodeURIComponent(craigslistSearchURL),
            columnIndex: index
        };

        fetch("http://localhost:8080/api/" , {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(myJsonRequestObj)
        })
        .then(response => response.json())
        .then(data =>
            this.setState({
                users: data,
                queryResponseData: data.response,
                message: "ok!!",
                dummystate: "ok finally changed state up here!!",
          })
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }))
        ;
    }

    render() {
        console.log("render CraigslistQueryColumn: " + JSON.stringify(this.props))
        return (
            <div>
                <div className="mui--text-left">
                    <CraigslistQueryColumnUserInput 
                        doRequest={this.doRequest}
                        doDeleteRequest={this.props.doDeleteRequest}
                        url={this.props.url}
                        columnIndex={this.props.columnIndex}
                        />
                </div>
                <CraigslistQueryColumnResults 
                    results={this.state.queryResponseData}
                    hello="hello" 
                    myprop="<b>zwatef</b>"
                />
            </div>
        );
    }
}