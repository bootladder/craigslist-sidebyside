class CitySelector extends React.Component {
    render() {
        console.log("render CraigslistQueryColumn: " + JSON.stringify(this.props))
        return (
            <div>
                <select id="areaAbb" class="js-only">
                    <option value="bham">birmingham, AL</option>
                        <option value="albanyga">albany, GA</option>
                        <option value="athensga">athens, GA</option>
                        <option value="atlanta">atlanta</option>
                        <option value="auburn">auburn</option>
                        <option value="chattanooga">chattanooga</option>
                        <option value="clarksville">clarksville, TN</option>
                        <option value="columbusga">columbus, GA</option>
                        <option value="cookeville">cookeville</option>
                        <option value="dothan">dothan, AL</option>
                        <option value="gadsden">gadsden</option>
                        <option value="hattiesburg">hattiesburg</option>
                        <option value="huntsville">huntsville</option>
                        <option value="jackson">jackson, MS</option>
                        <option value="jacksontn">jackson, TN</option>
                        <option value="macon">macon</option>
                        <option value="memphis">memphis</option>
                        <option value="meridian">meridian</option>
                        <option value="mobile">mobile, AL</option>
                        <option value="montgomery">montgomery</option>
                        <option value="nashville">nashville</option>
                        <option value="northmiss">north MS</option>
                        <option value="nwga">northwest GA</option>
                        <option value="okaloosa">okaloosa</option>
                        <option value="pensacola">pensacola</option>
                        <option value="shoals">the shoals</option>
                        <option value="tuscaloosa">tuscaloosa</option>
                </select>
            </div>
        );
    }
}