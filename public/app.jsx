class CraigslistQueryColumnResults extends React.Component {

  render() {
    const html = $.parseHTML( this.props.results ) //console.log(html.html)
    
    //const b = $(html).find(".result-row").length; console.log(b)
    //const z = $(html).find(".result-row").html()
      
    var resultRows = []
    var i=0
    $(html).find(".result-row").each( function(q) {
        //console.log(this.innerHTML)
        i++
        resultRows.push(React.createElement(
            'div', {key: i}, <div dangerouslySetInnerHTML= {{__html: this.innerHTML}} />
        ))
    })

    var stations = [];
    var i
    for(i=0; i<4; i++){
      stations.push(
          <div key={i} className="station">
              Call: {i}, Freq: {i} <br/>
          </div>
      )
    }

    return (
      <div>
            <div className="stations">{stations}</div>
            <div className="hello">{resultRows}</div>
      </div>
    );
  }
}

class CraigslistQueryColumnUserInput extends React.Component {
  render() {
      return (
          <div>
              URL: <input type="text" name="url"/><br/>
              <button onClick={this.props.buttonClicked} /> <br />

              <input value={37} onChange={this.props.handleChange} /> <br />
              hello ocomponent
          </div>
      );
  }
}

class CraigslistQueryColumn extends React.Component {
    state = {
        message: "initial state",
        users: [],
        error: null
    }
    constructor(props) {
      super(props);
      this.handleChange = this.handleChange.bind(this);
    }
    handleChange(e) {
        this.props.onInputEvent(e.target.value);
    }

    buttonClicked(e) {
        console.log("button clicked" + e)
    }

    componentDidMount() {

        const myURL="https://detroit.craigslist.org/d/architect-engineer-cad/search/egr"

        var myObj = {
            searchURL: encodeURIComponent(myURL),
        };

        fetch("http://localhost:8080/api/" , {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(myObj)
        })
        .then(response => response.json())
        .then(data =>
            this.setState({
                users: data,
                craigslistQueryResponse: data.response,
                message: "ok!!",
          })
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }

    render() {
        return (
            <div>
                <CraigslistQueryColumnUserInput 
                    buttonClicked={this.buttonClicked}
                    handleChange={this.handleChange}/>
                <div> 
                    The Message is: {JSON.stringify(this.state.message)}
                </div>
                <CraigslistQueryColumnResults 
                      results={this.state.craigslistQueryResponse}
                      hello="hello" 
                      myprop="<b>zwatef</b>"
                />
            </div>
        );
    }
}

class App extends React.Component {
    handleQueryInput(e) {
        console.log("hello"+e)
          }
          render() {
              return (
                  <div>
                      <h1>This is the App</h1>
                      <hr/>
                      <div className="search-table-outter wrapper">
                          <table className="search-table inner">
                          <tbody>
                            <tr>
                              <td>
                                <div>
                                  craigslist-fetch/result/result1
                                    <CraigslistQueryColumn 
                                        onInputEvent={this.handleQueryInput}
                                    />
                                </div>
                              </td> 
                              <td>
                                <div>
                                  craigslist-fetch/result/result2
                                    <CraigslistQueryColumn 
                                        onInputEvent={this.handleQueryInput}
                                    />
                                </div>
                              </td>
                              <td>
                                <div>
                                  craigslist-fetch/result/result2
                                    <CraigslistQueryColumn 
                                        onInputEvent={this.handleQueryInput}
                                    />
                                </div>
                              </td>
                            </tr>
                          </tbody>
                          </table>
                      </div>


                <CraigslistQueryColumn onInputEvent={this.handleQueryInput}/>
            </div>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('app'));

