#include <opencv2/imgproc/imgproc.hpp>
#include <opencv2/highgui/highgui.hpp>
#include <iostream>
#include "cuda_runtime.h"
#include "device_launch_parameters.h"


using namespace std;
using namespace cv;

/// Global Variables





/// Function headers

//void BWFilter(Mat image);
void JustRed(unsigned char* Input_Image, int Width, int Height, int Channels, string filterType);

__global__ void red(unsigned char* image);
__global__ void BWFilter(unsigned char* Image, int Channels);
__global__ void SepiaFilter(unsigned char* Image, int Channels);
__global__ void AvatarFilter(unsigned char* Image, int Channels);

/**
 * function main
 */

int main(int argc, char** argv)
{
    Mat src; Mat dst;
    /// Load the source image
    src = imread(argv[1]);
    string name = argv[1];
    string filterType = argv[2];
    
    dst = src.clone();

    
    JustRed(dst.data, src.cols, src.rows, src.channels() , filterType);
    
    string output_name = "../public/results/"+name;
    imwrite(output_name,dst);
    
    return 0;
}

/*
void BWFilter(Mat image){

    for (int i = 0; i < image.cols; i++){
        for (int j = 0; j < image.rows; j++){
            Vec3b colors = image.at<cv::Vec3b>(j,i);
            int blue = colors[0];
            int green = colors[1];
            int red = colors[2];
            colors[0] = colors[1] = colors[2] = (red + green + blue) /3;
            image.at<cv::Vec3b>(j,i) = colors;
        }
    }


    imshow(window_name,image);
    waitKey(0);
    return;

}
*/
void JustRed(unsigned char* Input_Image, int Width, int Height, int Channels, string filterType){

    
    unsigned char* Dev_Input_Image = NULL;
   
     //allocate the memory in gpu
     cudaMalloc((void**)&Dev_Input_Image, Height * Width * Channels);
    
     //copy data from CPU to GPU
     cudaMemcpy(Dev_Input_Image, Input_Image, Height * Width * Channels, cudaMemcpyHostToDevice);
 
     dim3 Grid_Image(Width, Height);
    
     if(filterType == "bw"){
        BWFilter<<<Grid_Image, 8 >>>(Dev_Input_Image, Channels);
     }
     else if(filterType == "sepia"){
        SepiaFilter<<<Grid_Image, 8 >>>(Dev_Input_Image, Channels);
     }
     else if(filterType == "avatar"){
        AvatarFilter<<<Grid_Image, 8 >>>(Dev_Input_Image, Channels);
     }
     
     //copy processed data back to cpu from gpu
     cudaMemcpy(Input_Image, Dev_Input_Image, Height * Width * Channels, cudaMemcpyDeviceToHost);
     
     //free gpu mempry
     cudaFree(Dev_Input_Image);
    
    return;

}


__global__ void SepiaFilter(unsigned char* Image, int Channels) {
    
    int x = blockIdx.x;
    int y = blockIdx.y;
    int idx = (x + y * gridDim.x) * Channels;

    if(Channels == 3){
        int B = Image[idx + 0];
        int G = Image[idx + 1];
        int R = Image[idx + 2];
    
        float tr = 0.393*R + 0.769*G + 0.189*B;
        float tg = 0.349*R + 0.686*G + 0.168*B;
        float tb = 0.272*R + 0.534*G + 0.131*B;

        if(tr > 255){
            tr = 255;
        }
        if(tg > 255){
            tg = 255;
        }
        if(tb > 255){
            tb = 255;
        }
        
        Image[idx + 0] = tb;
        Image[idx + 1] = tg;
        Image[idx + 2] = tr;

    }

}

__global__ void BWFilter(unsigned char* Image, int Channels) {
    int x = blockIdx.x;
    int y = blockIdx.y;
    int idx = (x + y * gridDim.x) * Channels;

    int suma = 0;
    for (int i = 0; i < Channels; i++) {
        
        suma = suma + Image[idx + i];
       
    }

    suma = suma / Channels;

    for (int i = 0; i < Channels; i++) {
        
        Image[idx + i] = suma;
       
    }
    
}


__global__ void AvatarFilter(unsigned char* Image, int Channels) {
    int x = blockIdx.x;
    int y = blockIdx.y;
    int idx = (x + y * gridDim.x) * Channels;

    if(Channels == 3){
        int B = Image[idx + 0];
        int G = Image[idx + 1];
        int R = Image[idx + 2];

        Image[idx + 0] = R;
        Image[idx + 1] = G;
        Image[idx + 2] = B;

    }
    
}


