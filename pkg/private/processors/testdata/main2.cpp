#include<fstream>

using namespace std;

int main() {
    ifstream fin("b.in");
    ofstream fout("b.out");
    int a = 0, b = 0;
    fin >> a >> b;
    fout << "a + b (file) = " << a + b << endl;
    return 0;
}